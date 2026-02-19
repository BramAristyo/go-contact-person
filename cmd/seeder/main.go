package main

import (
	"context"
	"log"
	"math/rand"

	"github.com/BramAristyo/rest-api-contact-person/internal/config"
	"github.com/BramAristyo/rest-api-contact-person/internal/database"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	db := database.Connect(cfg.DatabaseUrl)
	defer db.Close()

	ctx := context.Background()
	SeedContacts(ctx, db, 5000)
	SeedGroups(ctx, db)
	SeedContactGroups(ctx, db)
}

func SeedContacts(ctx context.Context, db *pgxpool.Pool, total int) {
	batch := &pgx.Batch{}

	for i := 0; i < total; i++ {
		name := faker.Name()
		email := faker.Email()
		phone := faker.Phonenumber()

		batch.Queue(
			`INSERT INTO contacts (name, email, phone) 
             VALUES ($1, $2, $3) 
             ON CONFLICT (email) DO NOTHING`,
			name, email, phone,
		)
	}

	br := db.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			log.Fatal("Batch Closed error")
		}
	}(br)

	for i := 0; i < total; i++ {
		_, err := br.Exec()
		if err != nil {
			log.Printf("skip row %d: %v\n", i+1, err)
		}
	}

	log.Printf("%d contacts seeded\n", total)
}
func SeedGroups(ctx context.Context, db *pgxpool.Pool) {
	groups := []string{
		"Family",
		"Friends",
		"Work",
		"Colleagues",
		"School",
		"University",
		"Clients",
		"Business",
		"Emergency",
		"Neighbors",
		"Gym",
		"Sports",
		"Hobbies",
		"Travel",
		"Medical",
		"Services",
		"Favorites",
		"Blocked",
	}

	batch := &pgx.Batch{}
	for _, g := range groups {
		batch.Queue(
			`INSERT INTO groups (name) VALUES ($1) ON CONFLICT DO NOTHING`, g,
		)
	}

	br := db.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {

		}
	}(br)

	for range groups {
		if _, err := br.Exec(); err != nil {
			log.Printf("skip group: %v \n", err)
		}
	}

	log.Printf("%d groups seeded\n", len(groups))
}
func SeedContactGroups(ctx context.Context, db *pgxpool.Pool) {
	// Streaming approach to avoid loading all contact and group ids into memory at once
	contactRows, err := db.Query(ctx, `SELECT id FROM contacts`)
	if err != nil {
		log.Fatal("Error while fetching contacts.")
	}
	defer contactRows.Close()

	var contactIds []int64
	for contactRows.Next() {
		var id int64
		if err := contactRows.Scan(&id); err != nil {
			return
		}
		contactIds = append(contactIds, id)
	}

	groupRows, err := db.Query(ctx, `SELECT id FROM groups`)
	if err != nil {
		log.Fatal("Error while fetching groups")
	}
	defer groupRows.Close()

	var groupIds []int64
	for groupRows.Next() {
		var id int64
		if err := groupRows.Scan(&id); err != nil {
			return
		}
		groupIds = append(groupIds, id)
	}

	batch := &pgx.Batch{}
	total := 0

	for _, contactId := range contactIds {
		rand.Shuffle(len(groupIds), func(i, j int) {
			groupIds[i], groupIds[j] = groupIds[j], groupIds[i]
		})

		for _, groupId := range groupIds[:3] {
			batch.Queue(
				`INSERT INTO contact_groups (contact_id, group_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`,
				contactId, groupId,
			)
			total++
		}
	}

	br := db.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {

		}
	}(br)

	for i := 0; i < total; i++ {
		if _, err := br.Exec(); err != nil {
			log.Printf("skip row %d: %v", i+1, err)
		}

	}

	log.Printf("%d contact_groups seeded\n\n", total)
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go func() { // Proses pesan dalam goroutine untuk tidak menghalangi event lain
			if v.Message.GetConversation() != "" {
				fmt.Println("Pesan diterima:", v.Message.GetConversation())
				time.Sleep(10 * time.Millisecond) // Tambahkan delay untuk mengurangi beban CPU
			}
		}()
	}
}

func main() {
	// Konfigurasi log database dan client
	dbLog := waLog.Stdout("Database", "WARN", true)
	clientLog := waLog.Stdout("Client", "WARN", true)

	// Membuat database store untuk menyimpan session
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	// Inisialisasi WhatsMeow client
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// Proses login dan koneksi
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(err)
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				fmt.Println("QR Code untuk login:", evt.Code)
			} else if evt.Event == "success" {
				fmt.Println("Login berhasil!")
				break
			}
		}
	} else {
		// Jika sudah login, langsung koneksi
		err = client.Connect()
		if err != nil {
			panic(err)
		}
	}

	// Menangkap sinyal Ctrl+C atau SIGTERM untuk menutup aplikasi dengan aman
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Menutup koneksi...")
	client.Disconnect()
	fmt.Println("Selesai.")
}

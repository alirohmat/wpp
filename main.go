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

	// Import driver SQLite modernc
	_ "modernc.org/sqlite"
)

func eventHandler(evt interface{}) {
	switch v := evt.(type) {
	case *events.Message:
		go func() { // Proses pesan dalam goroutine
			if v.Message.GetConversation() != "" {
				fmt.Println("Pesan diterima:", v.Message.GetConversation())
				time.Sleep(10 * time.Millisecond) // Tambahkan delay kecil untuk mengurangi beban CPU
			}
		}()
	}
}

func main() {
	// Konfigurasi log untuk database dan client
	dbLog := waLog.Stdout("Database", "WARN", true)
	clientLog := waLog.Stdout("Client", "WARN", true)

	// Membuat database store menggunakan driver modernc SQLite
	container, err := sqlstore.New("sqlite", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(fmt.Sprintf("failed to open database: %v", err))
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		panic(fmt.Sprintf("failed to get first device: %v", err))
	}

	// Inisialisasi client WhatsMeow
	client := whatsmeow.NewClient(deviceStore, clientLog)
	client.AddEventHandler(eventHandler)

	// Login atau koneksi ulang
	if client.Store.ID == nil {
		qrChan, _ := client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			panic(fmt.Sprintf("failed to connect: %v", err))
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
		// Koneksi ulang
		err = client.Connect()
		if err != nil {
			panic(fmt.Sprintf("failed to connect: %v", err))
		}
	}

	// Menangkap sinyal untuk menghentikan aplikasi
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("Menutup koneksi...")
	client.Disconnect()
	fmt.Println("Selesai.")
}

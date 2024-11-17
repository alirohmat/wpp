package main

import (
    "errors"
    "fmt"
    "log"

    "github.com/tulir/whatsmeow"
    "github.com/tulir/whatsmeow/events"
    "github.com/tulir/whatsmeow/store/sqlstore"
)

func main() {
    // Inisialisasi penyimpanan device
    deviceStore, err := sqlstore.New("sqlite3", "file:example.db?_foreign_keys=on", nil)
    if err!= nil {
        log.Fatal("Gagal inisialisasi penyimpanan device:", err)
    }

    // Buat klien WhatsApp baru
    client := whatsmeow.NewClient(deviceStore, nil)

    // Pastikan klien terhubung sebelum melakukan pairing
    if err := client.Connect(); err!= nil {
        log.Fatal("Gagal menghubungkan klien:", err)
    }

    // Tunggu hingga koneksi fully established (menggunakan whatsmeow/events)
    eventoChan := make(chan whatsmeow.Event)
    eventHandlerID := client.AddEventHandler(func(event whatsmeow.Event) {
        eventoChan <- event
    })
    go func() {
        for {
            select {
            case event := <-eventoChan:
                if _, ok := event.(*events.QRCodeEvent); ok {
                    // Koneksi sekarang sudah fully established
                    close(eventoChan)
                    return
                }
            }
        }
    }()
    // Tunggu hingga event QRCodeEvent diterima
    <-eventoChan

    // Contoh penggunaan PairPhone()
    phone := "6287834100533" // Nomor telepon yang ingin dipasangkan
    showPushNotification := true
    clientType := whatsmeow.PairClientOtherWebClient
    clientDisplayName := "Browser (Windows)" // Format wajib: Browser (OS)
    code, err := client.PairPhone(phone, showPushNotification, clientType, clientDisplayName)

    if err!= nil {
        switch {
        case errors.Is(err, &whatsmeow.PairDatabaseError{}):
            log.Println("Kesalahan Basis Data Pairing:", err)
        case errors.Is(err, &whatsmeow.PairProtoError{}):
            log.Println("Kesalahan Protokol Pairing:", err)
        default:
            log.Println("Kesalahan tidak diketahui pada PairPhone():", err)
        }
    } else {
        fmt.Printf("Kode Pairing untuk %s: %s\n", phone, code)
    }

    // Jangan lupa untuk menghapus event handler jika tidak dibutuhkan lagi
    client.RemoveEventHandler(eventHandlerID)
}

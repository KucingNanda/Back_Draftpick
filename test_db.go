package main

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Karakter struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Nama       string `gorm:"type:varchar(100);not null" json:"nama"`
	Elemen     string `gorm:"type:varchar(50)" json:"elemen"`
	Senjata    string `gorm:"type:varchar(50)" json:"senjata"`
	Rarity     string `gorm:"type:varchar(50)" json:"rarity"`
	Cost       int    `gorm:"type:int" json:"cost"`
	ImageURL   string `gorm:"type:text" json:"imageUrl"`
	Konstelasi string `gorm:"type:varchar(100)" json:"konstelasi"`
	UlangTahun string `gorm:"type:varchar(50)" json:"ulangTahun"`
	Region     string `gorm:"type:varchar(50);default:'Tidak Ada'" json:"region"`
	Sinergi    string `gorm:"type:varchar(50);default:'Tidak Ada'" json:"sinergi"`
}

func main() {
	dsn := "host=postgresql-kucing27.alwaysdata.net user=kucing27 password=septian27 dbname=kucing27_draftpick port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Gagal konek DB: ", err)
	}

	fmt.Println("Koneksi berhasil, mencoba insert...")

	k := Karakter{
		Nama: "TestCharacter",
		Elemen: "Pyro",
		Senjata: "Sword",
		Rarity: "4 Star",
		Cost: 0,
		ImageURL: "test",
		Konstelasi: "Test",
		UlangTahun: "1 Jan",
		Region: "Mondstadt",
		Sinergi: "Tidak Ada",
	}

	result := db.Create(&k)
	if result.Error != nil {
		fmt.Println("ERROR INSERT:", result.Error)
	} else {
		fmt.Println("BERHASIL INSERT! ID:", k.ID)
		// cleanup
		db.Delete(&k)
	}
}

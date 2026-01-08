package main
import (
	"fmt"
	"log"
	"github.com/Czeeen/lets-go-programming-v2025-autumn-spbstu-40008/prokopenko.nikita/task-8/internal/config"
)
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Printf("%s %s\n", cfg.AppStatus, cfg.ReportLevel)
}

package spi

type Service interface {
	Close() error
}
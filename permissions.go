package main

type KPermission struct {
	Name string `yaml:"name"`
	Slug string `yaml:"slug"`
}

type KRole struct {
	ID          string   `yaml:"id"`
	Permissions []string `yaml:"permissions"`
}

type KMember struct {
	ID          string   `yaml:"id"`
	Permissions []string `yaml:"permissions"`
	Roles       []string `yaml:"roles"`
}

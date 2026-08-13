data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "-mod=mod",
    "./cmd/atlas/atlas.go"
  ]
}

env "dev" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/16/dev?search_path=public"
  
  migration {
    dir = "file://migrations?format=goose"
  }

  format {
    migrate {
      diff = "{{ sql . \"  \"}}"
    }
  }
}

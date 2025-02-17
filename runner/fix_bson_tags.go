package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

func main() {
	filePath := "graph/model/models_gen.go" // Ajusta la ruta según tu estructura

	// Leer el archivo
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error al leer el archivo:", err)
		return
	}

	content := string(data)

	// Reemplazar json:"id" por json:"id" bson:"_id,omitempty"
	idRegex := regexp.MustCompile(`json:"id"`)
	content = idRegex.ReplaceAllString(content, `json:"id" bson:"_id,omitempty"`)

	// Expresión regular para encontrar otros campos con etiquetas JSON
	fieldRegex := regexp.MustCompile(`json:"(\w+)"`)
	content = fieldRegex.ReplaceAllStringFunc(content, func(match string) string {
		fieldName := fieldRegex.FindStringSubmatch(match)[1]
		if fieldName == "id" {
			return match // Ya fue modificado antes
		}
		return fmt.Sprintf(`json:"%s" bson:"%s,omitempty"`, fieldName, fieldName)
	})

	// Escribir los cambios en el archivo
	err = ioutil.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		fmt.Println("Error al escribir el archivo:", err)
		return
	}

	fmt.Println("✅ Archivo actualizado correctamente.")
}

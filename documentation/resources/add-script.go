/*
  COLMENA-DESCRIPTION-SERVICE
  Copyright © 2024 EVIDEN
 
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
 
  http://www.apache.org/licenses/LICENSE-2.0
 
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
 
  This work has been implemented within the context of COLMENA project.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const comment = `/*
  COLMENA-DESCRIPTION-SERVICE
  Copyright © 2024 EVIDEN
 
  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at
 
  http://www.apache.org/licenses/LICENSE-2.0
 
  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
 
  This work has been implemented within the context of COLMENA project.
*/
`

func main() {
	root := "." // Directorio raíz del proyecto
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			err = addCommentToFile(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking the path %q: %v\n", root, err)
	}
}

func addCommentToFile(path string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(input), "\n")

	// Comprobar si el comentario ya está presente
	if strings.HasPrefix(lines[0], "// Este es el comentario añadido a todos los archivos Go") {
		return nil
	}

	// Añadir el comentario al inicio
	output := comment + strings.Join(lines, "\n")
	err = os.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

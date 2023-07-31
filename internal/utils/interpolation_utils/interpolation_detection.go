package interpolation_utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/aymerick/raymond"
)

func DetectInterpolation(template string, vars map[string]string) error {
	output, err := raymond.Parse(template)
	if err != nil {
		return fmt.Errorf("Error parsing template: %s", template)
	}
	variableNames := make(map[string]bool) //using a map as a set

	ast := strings.Split(output.PrintAST(), "\n")
	pathRex := regexp.MustCompile(`\:\w+`) // results of the handblebars ast that contains variables are in the form of "{{ PATH:<Name> [] }}"
	for _, line := range ast {
		if strings.Contains(line, "PATH") {
			varName := pathRex.FindString(line)
			varName = strings.TrimPrefix(varName, ":")
			variableNames[varName] = true
			if _, ok := vars[varName]; !ok {
				return fmt.Errorf("Variable \"%s\" is not found in variables", varName)
			}
		}
	}

	return nil
}

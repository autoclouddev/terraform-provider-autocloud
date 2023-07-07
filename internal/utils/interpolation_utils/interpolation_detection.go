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
	variableNames := map[string]struct{}{} //using a map as a set

	ast := strings.Split(output.PrintAST(), "\n")
	pathRex := regexp.MustCompile(`\:\w+`) // results of the handblebars ast that contains variables are in the form of "{{ PATH:<Name> [] }}"
	for _, line := range ast {
		if strings.Contains(line, "PATH") {
			varName := pathRex.FindString(line)
			varName = strings.TrimPrefix(varName, ":")
			variableNames[varName] = struct{}{}
		}
	}
	// will accumulate error messages, to display all of them at once
	errorMsgs := make([]string, 0)
	for key := range vars {
		if _, ok := variableNames[key]; !ok {
			errorMsgs = append(errorMsgs, fmt.Sprintf("Variable \"%s\" is not found in template", key))
		} else {
			delete(variableNames, key)
		}
	}

	for key := range variableNames {
		errorMsgs = append(errorMsgs, fmt.Sprintf("Variable \"%s\" is not found in variables", key))
	}

	if len(errorMsgs) > 0 {
		return fmt.Errorf(strings.Join(errorMsgs, "\n"))
	}

	return nil
}

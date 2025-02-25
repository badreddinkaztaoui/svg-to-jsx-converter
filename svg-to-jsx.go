package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	inputFile := flag.String("input", "", "Input SVG file path")
	outputFile := flag.String("output", "", "Output TSX file path")
	componentName := flag.String("name", "SvgIcon", "React component name")
	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: Please provide an input file with -input flag")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *outputFile == "" {
		baseName := filepath.Base(*inputFile)
		nameWithoutExt := strings.TrimSuffix(baseName, filepath.Ext(baseName))
		*outputFile = nameWithoutExt + ".tsx"
	}

	svgContent, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		os.Exit(1)
	}

	tsxContent := convertSvgToTsx(string(svgContent), *componentName)

	err = ioutil.WriteFile(*outputFile, []byte(tsxContent), 0644)
	if err != nil {
		fmt.Printf("Error writing output file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully converted %s to %s\n", *inputFile, *outputFile)
}

func convertSvgToTsx(svgContent, componentName string) string {
	re := regexp.MustCompile(`<svg[^>]*>([\s\S]*?)</svg>`)
	svgMatch := re.FindStringSubmatch(svgContent)
	
	var svgBody string
	if len(svgMatch) > 1 {
		svgBody = svgMatch[1]
	}

	attrRe := regexp.MustCompile(`<svg([^>]*)>`)
	attrMatch := attrRe.FindStringSubmatch(svgContent)
	var svgAttrs string
	if len(attrMatch) > 1 {
		svgAttrs = attrMatch[1]
	}

	attrs := parseAttributes(svgAttrs)

	var b strings.Builder

	b.WriteString("import * as React from \"react\";\n\n")
	b.WriteString(fmt.Sprintf("const %s: React.FC<React.SVGProps<SVGElement>> = (props) => (\n", componentName))
	b.WriteString("  <svg\n")
	b.WriteString("    xmlns=\"http://www.w3.org/2000/svg\"\n")

	for key, value := range attrs {
		if key != "xmlns" {
			b.WriteString(fmt.Sprintf("    %s=%s\n", key, value))
		}
	}

	b.WriteString("    {...props}\n")
	b.WriteString("  >\n")

	body := cleanSvgBody(svgBody)
	
	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) != "" {
			b.WriteString("    " + line + "\n")
		}
	}

	b.WriteString("  </svg>\n")
	b.WriteString(");\n\n")
	b.WriteString(fmt.Sprintf("export default React.memo(%s);\n", componentName))

	return b.String()
}

func parseAttributes(attrsStr string) map[string]string {
	attrs := make(map[string]string)
	
	re := regexp.MustCompile(`([a-zA-Z0-9_:-]+)\s*=\s*("[^"]*"|'[^']*')`)
	matches := re.FindAllStringSubmatch(attrsStr, -1)
	
	svgToReactAttrMap := map[string]string{
		"class":             "className",
		"clip-path":         "clipPath",
		"clip-rule":         "clipRule",
		"fill-opacity":      "fillOpacity",
		"fill-rule":         "fillRule",
		"stroke-dasharray":  "strokeDasharray",
		"stroke-dashoffset": "strokeDashoffset",
		"stroke-linecap":    "strokeLinecap",
		"stroke-linejoin":   "strokeLinejoin",
		"stroke-miterlimit": "strokeMiterlimit",
		"stroke-opacity":    "strokeOpacity",
		"stroke-width":      "strokeWidth",
		"text-anchor":       "textAnchor",
		"font-family":       "fontFamily",
		"font-size":         "fontSize",
		"font-weight":       "fontWeight",
		"xlink:href":        "xlinkHref",
		"xml:space":         "xmlSpace",
		"stop-color":        "stopColor",
		"stop-opacity":      "stopOpacity",
		"color-interpolation": "colorInterpolation",
		"color-rendering":   "colorRendering",
		"enable-background": "enableBackground",
		"dominant-baseline": "dominantBaseline",
		"shape-rendering":   "shapeRendering",
		"text-decoration":   "textDecoration",
		"vector-effect":     "vectorEffect",
	}
	
	for _, match := range matches {
		if len(match) >= 3 {
			name := match[1]
			value := match[2]
			
			if mappedName, exists := svgToReactAttrMap[name]; exists {
				name = mappedName
			} else if name != "viewBox" && strings.Contains(name, "-") {
				parts := strings.Split(name, "-")
				name = parts[0]
				for i := 1; i < len(parts); i++ {
					if len(parts[i]) > 0 {
						name += strings.ToUpper(parts[i][:1]) + parts[i][1:]
					}
				}
			}
			
			attrs[name] = value
		}
	}
	
	return attrs
}

func cleanSvgBody(body string) string {
	lines := strings.Split(body, "\n")
	var result []string
	
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			trimmed = convertAttributesInLine(trimmed)
			result = append(result, trimmed)
		}
	}
	
	return strings.Join(result, "\n")
}

func convertAttributesInLine(line string) string {
	svgToReactAttrMap := map[string]string{
		"clip-path":         "clipPath",
		"clip-rule":         "clipRule",
		"fill-opacity":      "fillOpacity",
		"fill-rule":         "fillRule",
		"stroke-dasharray":  "strokeDasharray",
		"stroke-dashoffset": "strokeDashoffset",
		"stroke-linecap":    "strokeLinecap",
		"stroke-linejoin":   "strokeLinejoin",
		"stroke-miterlimit": "strokeMiterlimit",
		"stroke-opacity":    "strokeOpacity",
		"stroke-width":      "strokeWidth",
		"text-anchor":       "textAnchor",
		"font-family":       "fontFamily",
		"font-size":         "fontSize",
		"font-weight":       "fontWeight",
		"stop-color":        "stopColor",
		"stop-opacity":      "stopOpacity",
		"enable-background": "enableBackground",
		"class":             "className",
	}
	
	for svgAttr, reactAttr := range svgToReactAttrMap {
		re := regexp.MustCompile(`(^|\s)` + svgAttr + `=`)
		line = re.ReplaceAllString(line, "${1}"+reactAttr+"=")
	}
	
	attrRe := regexp.MustCompile(`\s([a-z]+-[a-z0-9-]+)=`)
	matches := attrRe.FindAllStringSubmatch(line, -1)
	
	for _, match := range matches {
		if len(match) >= 2 {
			kebabAttr := match[1]
			alreadyHandled := false
			for svgAttr := range svgToReactAttrMap {
				if svgAttr == kebabAttr {
					alreadyHandled = true
					break
				}
			}
			
			if !alreadyHandled {
				camelAttr := convertKebabToCamel(kebabAttr)
				line = strings.ReplaceAll(line, " "+kebabAttr+"=", " "+camelAttr+"=")
			}
		}
	}
	
	return line
}

func convertKebabToCamel(kebab string) string {
	parts := strings.Split(kebab, "-")
	result := parts[0]
	
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			result += strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	
	return result
}

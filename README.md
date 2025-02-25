# SVG to TSX Converter

A simple command-line tool written in Go that converts SVG files into React TypeScript components (TSX).

## Features

- Converts SVG files to React functional components
- Automatically transforms SVG attributes to React camelCase format
- Preserves viewBox and other important attributes
- Adds proper TypeScript types using React.SVGProps
- Memoizes the component with React.memo for performance
- Allows custom component naming

## Installation

### Prerequisites

- Go 1.16 or higher

### Building from Source

Clone this repository and build the executable:

```bash
go build -o svg-to-tsx
```

## Usage

```bash
./svg-to-tsx -input <svg-file> -output <tsx-file> -name <component-name>
```

### Command-line Arguments

- `-input` (required): Path to the input SVG file
- `-output` (optional): Path to the output TSX file (defaults to the input filename with .tsx extension)
- `-name` (optional): React component name (defaults to "SvgIcon")

### Example

```bash
./svg-to-tsx -input icons/menu.svg -output components/MenuIcon.tsx -name MenuIcon
```

## How It Works

The converter:

1. Reads the SVG file content
2. Extracts the SVG attributes and body
3. Converts kebab-case attributes (e.g., `fill-rule`) to camelCase (`fillRule`)
4. Generates a properly formatted React functional component
5. Writes the result to the specified output file

## Output Format

The generated TSX file follows this structure:

```tsx
import * as React from "react";

const ComponentName: React.FC<React.SVGProps<SVGElement>> = (props) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    width="24"
    height="24"
    {...props}
  >
    {/* SVG content here */}
  </svg>
);

export default React.memo(ComponentName);
```

## Supported SVG Attribute Conversions

The script handles common SVG attributes and converts them to their React equivalents:

- `class` → `className`
- `clip-path` → `clipPath`
- `fill-rule` → `fillRule`
- `stroke-width` → `strokeWidth`
- And many more...

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## Limitations

- The converter assumes well-formed SVG files
- Complex SVG files with embedded scripts or styles may require additional handling

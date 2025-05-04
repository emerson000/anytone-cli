# Anytone CLI

A proof-of-concept command-line interface (CLI) for working with Anytone codeplugs. This tool allows you to view and modify parameters in Anytone radio codeplug (.rdt) files without using the official CPS software.

> [!WARNING]
> This is a work-in-progress and not suitable for production use. Contributions are welcome. Use at your own risk.

## Features

- View codeplug information (model, radio IDs, etc.)
- Update radio IDs (DMR ID numbers)
- Works with AT-D878UVII

## Installation

### Prerequisites

- Python 3.6 or higher

### From PyPI

```bash
pip install anytone-cli
```

### From Source

```bash
git clone https://github.com/emerson000/anytone-cli.git
cd anytone-cli
pip install -e .
```

## Usage

The CLI requires a codeplug file (.rdt) as its first argument, followed by a command:

```bash
anytone-cli <codeplug_file.rdt> <command> [options]
```

### Commands

#### View Codeplug Information

```bash
anytone-cli codeplug.rdt info
```

This displays general information about your codeplug file, including:
- Radio IDs configured

#### Update Radio ID

```bash
anytone-cli codeplug.rdt update radio_id <index> <new_id>
```

Example:
```bash
anytone-cli codeplug.rdt update radio_id 0 3161234
```

This updates the first radio ID (index 0) to 3161234.

### Examples

To display information about a codeplug:
```bash
anytone-cli my_radio.rdt info
```

To update the second radio ID:
```bash
anytone-cli my_radio.rdt update radio_id 1 3165678
```

## License

This project is licensed under the MIT License.


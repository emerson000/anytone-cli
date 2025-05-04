#!/usr/bin/env python3

import argparse
import sys
from anytone_cli.codeplug import display_info

def main():
    parser = argparse.ArgumentParser(
        description="Anytone CLI - Command line interface for Anytone radios",
        prog="anytone-cli"
    )
    
    # Create subparsers for different commands
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")
    
    # Add 'info' command
    info_parser = subparsers.add_parser("info", help="Display information about a codeplug")
    info_parser.add_argument("codeplug_file", help="Path to the codeplug file")
    
    # Parse the arguments
    args = parser.parse_args()
    
    # Execute the chosen command
    if args.command == "info":
        print("Executing info command...")
        display_info(args.codeplug_file)
    else:
        parser.print_help()
        return 1
    
    return 0

if __name__ == "__main__":
    sys.exit(main()) 
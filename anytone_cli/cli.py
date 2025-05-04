#!/usr/bin/env python3

import argparse
import sys
from anytone_cli.codeplug import display_info, update_radio_id

def main():
    parser = argparse.ArgumentParser(
        description="Anytone CLI - Command line interface for Anytone radios",
        prog="anytone-cli"
    )
    
    # Add the codeplug file as a positional argument before subcommands
    parser.add_argument("codeplug_file", help="Path to the codeplug file")
    
    # Create subparsers for different commands
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")
    
    # Add 'info' command
    info_parser = subparsers.add_parser("info", help="Display information about a codeplug")
    
    # Add 'update' command
    update_parser = subparsers.add_parser("update", help="Update elements in the codeplug")
    update_subparsers = update_parser.add_subparsers(dest="update_command", help="Element to update")
    
    # Add 'radio_id' subcommand to 'update'
    radio_id_parser = update_subparsers.add_parser("radio_id", help="Update a radio ID")
    radio_id_parser.add_argument("index", type=int, help="Index of the radio ID to update")
    radio_id_parser.add_argument("radio_id", type=int,help="New radio ID value")
    
    # Parse the arguments
    args = parser.parse_args()
    
    # Execute the chosen command
    if args.command == "info":
        print("Executing info command...")
        display_info(args.codeplug_file)
    elif args.command == "update":
        if args.update_command == "radio_id":
            update_radio_id(args.codeplug_file, args.index, args.radio_id)
        else:
            update_parser.print_help()
            return 1
    elif args.command is None:
        parser.print_help()
        return 1
    
    return 0

if __name__ == "__main__":
    sys.exit(main()) 
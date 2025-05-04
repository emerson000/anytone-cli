"""
Codeplug module for Anytone CLI

This module provides functionality for working with Anytone radio codeplugs.
"""

import os
import datetime
import sys
from anytone_cli.radio_ids import (
    find_radio_id_section, 
    extract_radio_ids, 
    update_radio_id as update_id_in_data, 
    RadioIdError
)


class CodeplugError(Exception):
    """Exception raised for errors in the codeplug handling."""
    pass


def read_codeplug(filename):
    """
    Read and parse an Anytone codeplug file.
    
    Args:
        filename (str): Path to the codeplug file
        
    Returns:
        dict: Dictionary containing codeplug information
        
    Raises:
        CodeplugError: If the file cannot be read or is not a valid codeplug
    """
    if not os.path.exists(filename):
        raise CodeplugError(f"File not found: {filename}")
    
    try:
        with open(filename, 'rb') as f:
            data = f.read()
            
        if len(data) < 16:
            raise CodeplugError("File too small to be a valid codeplug")

        codeplug_info = {
            'filename': filename,
            'filesize': len(data),
            'model': 'Unknown Anytone Model',
            'format_version': '1.0',
            'last_modified': datetime.datetime.fromtimestamp(os.path.getmtime(filename)),
            'radio_ids': [],
            'channels': 0,
            'zones': 0,
            'contacts': 0,
        }
        
        # Check specific model identifiers at known addresses
        if data[0:4] == b'D878':
            codeplug_info['model'] = 'Anytone AT-D878UV'
        elif data[0:4] == b'D578':
            codeplug_info['model'] = 'Anytone AT-D578UV'
        
        # Check for D878UVII at address 0x00000009 (9) through 0x00000010 (16)
        if len(data) >= 17:  # Ensure we have enough data (0-indexed, so need 17 bytes for index 16)
            model_bytes = data[9:17]
            if model_bytes == b'D878UVII':
                codeplug_info['model'] = 'Anytone AT-D878UVII'
                codeplug_info['model_bytes'] = model_bytes.hex()
        
        # Find the radio ID section
        try:
            offset, offset_end = find_radio_id_section(data)
            radio_id_data = data[offset:offset_end]
            
            # Extract radio IDs from the data section
            codeplug_info['radio_ids'] = extract_radio_ids(radio_id_data)
        except RadioIdError as e:
            print(f"Warning: {str(e)}", file=sys.stderr)
        
        return codeplug_info
        
    except IOError as e:
        raise CodeplugError(f"Error reading codeplug: {str(e)}")


def update_radio_id(codeplug_file, index, radio_id):
    """
    Update the radio ID at the given index with the new radio ID.
    
    Args:
        codeplug_file (str): Path to the codeplug file
        index (int): Index of the radio ID to update (1-based)
        radio_id (int): New radio ID value (must be a 24-bit value, 1-16777215)
        
    Raises:
        CodeplugError: If the file cannot be read/written or if the radio ID is invalid
    """
    if not os.path.exists(codeplug_file):
        raise CodeplugError(f"File not found: {codeplug_file}")
    
    try:
        # 1. Load the codeplug file
        with open(codeplug_file, 'rb') as f:
            data = bytearray(f.read())
        
        # Use the radio_ids module to update the radio ID
        try:
            modified_data = update_id_in_data(data, index, radio_id)
            
            # Save the modified codeplug file
            with open(codeplug_file, 'wb') as f:
                f.write(modified_data)
                
            print(f"Successfully updated radio ID at index {index} to {radio_id} in {codeplug_file}")
        except RadioIdError as e:
            raise CodeplugError(str(e))
            
    except IOError as e:
        raise CodeplugError(f"Error updating radio ID: {str(e)}")


def display_info(filename):
    """
    Display information about a codeplug file.
    
    Args:
        filename (str): Path to the codeplug file
    """
    try:
        info = read_codeplug(filename)
        
        print(f"\nCodeplug Information:")
        print(f"====================")
        print(f"Filename:       {info['filename']}")
        print(f"File size:      {info['filesize']:,} bytes")
        print(f"Radio model:    {info['model']}")
        if 'model_bytes' in info:
            print(f"Model bytes:    0x{info['model_bytes']} (at address 0x00000009-0x00000010)")
        print(f"Format version: {info['format_version']}")
        print(f"Last modified:  {info['last_modified'].strftime('%Y-%m-%d %H:%M:%S')}")
        print(f"Channels:       {info['channels']}")
        print(f"Zones:          {info['zones']}")
        print(f"Contacts:       {info['contacts']}")
        # Display radio IDs as a numbered list, one per line
        print(f"Radio IDs:")
        for i, radio_id in enumerate(info['radio_ids'], 1):
            print(f"  {i}. {radio_id}")
        
    except CodeplugError as e:
        print(f"Error: {str(e)}", file=sys.stderr) 
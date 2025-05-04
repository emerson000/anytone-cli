"""
Codeplug module for Anytone CLI

This module provides functionality for working with Anytone radio codeplugs.
"""

import os
import struct
import datetime
import sys


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
        
        pattern = bytes.fromhex("46 4F 20 42 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01")
        pattern_end = bytes.fromhex("00 01 00 04 00 00 01 00 02 00 03")
        offset = data.find(pattern) + len(pattern)
        offset_end = data.find(pattern_end)
        radio_id_data = data[offset:offset_end]
        
        # Extract radio IDs from the data section
        codeplug_info['radio_ids'] = _extract_radio_ids(radio_id_data)
        
        return codeplug_info
        
    except (IOError, struct.error) as e:
        raise CodeplugError(f"Error reading codeplug: {str(e)}")

def _extract_radio_ids(data):
    """
    Extract all radio IDs from the radio ID section data.
    
    Args:
        data (bytes): Radio ID section data
        
    Returns:
        list: List of radio IDs
    """
    radio_ids = []
    
    # Start at position 2 (first radio ID location)
    pos = 2
    
    # Process data until we run out of bytes
    while pos + 2 < len(data):
        # Extract the 24-bit radio ID
        radio_id = _get_uint24(data, pos)
        
        # Add valid IDs (greater than 0)
        if radio_id > 0:
            radio_ids.append(radio_id)
            
        # Look for the next potential ID location
        # Typically after a 0x00 byte that serves as a separator
        next_zero = data.find(b'\x00', pos + 3)
        if next_zero == -1:
            break
            
        # Move to the position after the separator
        pos = next_zero + 2
    
    return radio_ids

def _get_uint24(data, offset):
    return data[offset] | (data[offset+1] << 8) | (data[offset+2] << 16)

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
        print(f"Radio IDs:      {', '.join(str(id) for id in info['radio_ids'])}")
        
    except CodeplugError as e:
        print(f"Error: {str(e)}", file=sys.stderr) 
"""
Codeplug module for Anytone CLI

This module provides functionality for working with Anytone radio codeplugs.
"""

import os
import struct
import datetime
import sys

# Pattern constants for radio ID section identification
RADIO_ID_SECTION_START = bytes.fromhex("46 4F 20 42 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01")
RADIO_ID_SECTION_END = bytes.fromhex("00 01 00 04 00 00 01 00 02 00 03")


class CodeplugError(Exception):
    """Exception raised for errors in the codeplug handling."""
    pass


def _find_radio_id_section(data):
    """
    Locate the radio ID section in the codeplug data.
    
    Args:
        data (bytes or bytearray): Codeplug data
        
    Returns:
        tuple: (start_offset, end_offset) of the radio ID section
        
    Raises:
        CodeplugError: If the radio ID section cannot be located
    """
    section_start = data.find(RADIO_ID_SECTION_START) + len(RADIO_ID_SECTION_START)
    section_end = data.find(RADIO_ID_SECTION_END)
    
    if section_start <= len(RADIO_ID_SECTION_START) or section_end == -1:
        raise CodeplugError("Could not locate radio ID section in codeplug")
        
    return section_start, section_end


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
        
        # Find the radio ID section
        offset, offset_end = _find_radio_id_section(data)
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

def _set_uint24(value):
    """Convert an integer to a 3-byte (24-bit) sequence."""
    return bytes([
        value & 0xFF,
        (value >> 8) & 0xFF,
        (value >> 16) & 0xFF
    ])

def update_radio_id(codeplug_file, index, radio_id):
    """
    Update the radio ID at the given index with the new radio ID.
    
    Args:
        codeplug_file (str): Path to the codeplug file
        index (int): Index of the radio ID to update (0-based)
        radio_id (int): New radio ID value (must be a 24-bit value, 1-16777215)
        
    Raises:
        CodeplugError: If the file cannot be read/written or if the radio ID is invalid
    """
    if not os.path.exists(codeplug_file):
        raise CodeplugError(f"File not found: {codeplug_file}")
    
    if not isinstance(radio_id, int) or radio_id < 1 or radio_id > 0xFFFFFF:
        raise CodeplugError(f"Invalid radio ID: {radio_id}. Must be between 1 and 16777215.")
    
    try:
        # 1. Load the codeplug file
        with open(codeplug_file, 'rb') as f:
            data = bytearray(f.read())
        
        # Find the radio ID section
        section_start, section_end = _find_radio_id_section(data)
        radio_id_data = data[section_start:section_end]
        
        # Extract existing radio IDs
        radio_ids = _extract_radio_ids(radio_id_data)
        
        if index < 0 or index >= len(radio_ids):
            raise CodeplugError(f"Radio ID index {index} is out of range. Available indexes: 0-{len(radio_ids)-1}")
        
        # 2. Modify the radio ID at the specified index
        # We need to find the position of the radio ID in the data
        pos = section_start + 2  # Start at the first radio ID location
        found_index = 0
        
        while pos + 2 < len(data) and found_index < len(radio_ids):
            if found_index == index:
                # Found the index, replace the radio ID
                data[pos:pos+3] = _set_uint24(radio_id)
                break
                
            # Look for the next radio ID
            next_zero = data.find(b'\x00', pos + 3, section_end)
            if next_zero == -1:
                break
                
            pos = next_zero + 2
            found_index += 1
        
        # 3. Save the codeplug file
        with open(codeplug_file, 'wb') as f:
            f.write(data)
            
        print(f"Successfully updated radio ID at index {index} to {radio_id} in {codeplug_file}")
        
    except (IOError, struct.error) as e:
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
        print(f"Radio IDs:      {', '.join(str(id) for id in info['radio_ids'])}")
        
    except CodeplugError as e:
        print(f"Error: {str(e)}", file=sys.stderr) 
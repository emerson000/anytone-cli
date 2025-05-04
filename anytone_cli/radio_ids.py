"""
Radio ID module for Anytone CLI

This module provides functionality for working with Anytone radio IDs.
"""

import struct

# Pattern constants for radio ID section identification
RADIO_ID_SECTION_START = bytes.fromhex("46 4F 20 42 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 01")
RADIO_ID_SECTION_END = bytes.fromhex("00 01 00 04 00 00 01 00 02 00 03")


class RadioIdError(Exception):
    """Exception raised for errors in the radio ID handling."""
    pass


def find_radio_id_section(data):
    """
    Locate the radio ID section in the codeplug data.
    
    Args:
        data (bytes or bytearray): Codeplug data
        
    Returns:
        tuple: (start_offset, end_offset) of the radio ID section
        
    Raises:
        RadioIdError: If the radio ID section cannot be located
    """
    section_start = data.find(RADIO_ID_SECTION_START) + len(RADIO_ID_SECTION_START)
    section_end = data.find(RADIO_ID_SECTION_END)
    
    if section_start <= len(RADIO_ID_SECTION_START) or section_end == -1:
        raise RadioIdError("Could not locate radio ID section in codeplug")
        
    return section_start, section_end


def extract_radio_ids(data):
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
        radio_id = get_uint24(data, pos)
        
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


def get_uint24(data, offset):
    """Convert a 3-byte (24-bit) sequence to an integer."""
    return data[offset] | (data[offset+1] << 8) | (data[offset+2] << 16)


def set_uint24(value):
    """Convert an integer to a 3-byte (24-bit) sequence."""
    return bytes([
        value & 0xFF,
        (value >> 8) & 0xFF,
        (value >> 16) & 0xFF
    ])


def update_radio_id(data, index, radio_id):
    """
    Update the radio ID at the given index in the codeplug data.
    
    Args:
        data (bytearray): Codeplug data
        index (int): Index of the radio ID to update (1-based)
        radio_id (int): New radio ID value (must be a 24-bit value, 1-16777215)
        
    Returns:
        bytearray: Modified codeplug data
        
    Raises:
        RadioIdError: If the radio ID is invalid or index is out of range
    """
    if not isinstance(radio_id, int) or radio_id < 1 or radio_id > 0xFFFFFF:
        raise RadioIdError(f"Invalid radio ID: {radio_id}. Must be between 1 and 16777215.")
    
    try:
        # Find the radio ID section
        section_start, section_end = find_radio_id_section(data)
        radio_id_data = data[section_start:section_end]
        
        # Extract existing radio IDs
        radio_ids = extract_radio_ids(radio_id_data)
        
        # Convert 1-based index to 0-based for internal operations
        zero_based_index = index - 1
        
        if zero_based_index < 0 or zero_based_index >= len(radio_ids):
            raise RadioIdError(f"Radio ID index {index} is out of range. Available indexes: 1-{len(radio_ids)}")
        
        # Modify the radio ID at the specified index
        pos = section_start + 2  # Start at the first radio ID location
        found_index = 0
        
        while pos + 2 < len(data) and found_index < len(radio_ids):
            if found_index == zero_based_index:
                # Found the index, replace the radio ID
                data[pos:pos+3] = set_uint24(radio_id)
                break
                
            # Look for the next radio ID
            next_zero = data.find(b'\x00', pos + 3, section_end)
            if next_zero == -1:
                break
                
            pos = next_zero + 2
            found_index += 1
            
        return data
            
    except (struct.error) as e:
        raise RadioIdError(f"Error updating radio ID: {str(e)}") 
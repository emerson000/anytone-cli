#!/usr/bin/env python3
"""
Test script for Anytone CLI
"""

import subprocess
import sys

def run_command(cmd):
    """Run a command and print the output"""
    print(f"Running: {' '.join(cmd)}")
    result = subprocess.run(cmd, capture_output=True, text=True)
    print(f"Exit code: {result.returncode}")
    print("Output:")
    print(result.stdout)
    if result.stderr:
        print("Error:")
        print(result.stderr)
    print("-" * 50)
    return result

def main():
    """Run tests for the CLI"""
    # Test help
    run_command(["python", "-m", "anytone_cli.cli", "--help"])
    
    # Test info command
    run_command(["python", "-m", "anytone_cli.cli", "info"])
    
    # Test read command
    run_command(["python", "-m", "anytone_cli.cli", "read", "-o", "test_output.json", "-t", "channels"])
    
    # Test write command
    run_command(["python", "-m", "anytone_cli.cli", "write", "-i", "test_input.json", "-t", "contacts"])
    
    print("All tests completed")

if __name__ == "__main__":
    main() 
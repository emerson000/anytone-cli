from setuptools import setup, find_packages

setup(
    name="anytone-cli",
    version="0.1.0",
    description="Command line interface for Anytone radios",
    author="emerson000",
    author_email="me@emerson.sh",
    packages=find_packages(),
    entry_points={
        "console_scripts": [
            "anytone-cli=anytone_cli.cli:main",
        ],
    },
    install_requires=[
        "argparse",
    ],
    python_requires=">=3.6",
) 
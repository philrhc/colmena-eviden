from pathlib import Path
from setuptools import setup, find_packages

with Path('README.md').open() as file:
    README = file.read()

with Path('requirements.txt').open() as file:
    INSTALL_REQUIRES = file.readlines()

setup(
    name="metrics-etl",
    author='2024, EVIDEN BDS R&D',
    packages=find_packages(exclude=['tests']),
    maintainer="jorge.rubioc",
    maintainer_email="jorge.rubioc@eviden.com",
    long_description=README,
    install_requires=INSTALL_REQUIRES,
    long_description_content_type='text/markdown',
    license="OSI Approved :: Apache Software License",
    python_requires='>=3.10', 
    zip_safe=False,
    package_data={'': ['config/zenoh_config.json5']},
    include_package_data=True,
    entry_points={
        'console_scripts': [
            'metrics-etl=metrics_etl.cli:start_etl',
            ]},
)

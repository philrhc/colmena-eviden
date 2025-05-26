from .example_contextdata import CompanyPremises


__version__ = '0.0.0'
def main():
	device = None # Environment variable, JSON file, TBD.
	r = CompanyPremises().locate(device)

# Crawlie
A webcrawler to parse html docmumets

## Crawlie Functions

### Internal Links

Crawlie will go through all internal links to the website located withing the HTML.
It will return a list of all internal links, and how many times they are refenced.

### External Links

Crawlie will identify any external links, and return them as a list, and how many times they have been identified.

### Email Addresses

Crawlie will identify any email addresses within the HTML and return the them as a list.

### Phone Numbers

Crawlie will identify any phone numbers that are within the HTML and return them as a list.

#### Proposed features

- Identify any application contact details
- Identify any commercial products used for site construction (probs not, I think warpayzler does this fine)
- Identify and return server host location
- Identify reverse DNS results
- Identify other domains hosted on the same DNS server

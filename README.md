To help understand the flow of events. Legend:

Job related tasks will be sent to the topic abbreviated by JT
Job related results will be sent to the topic abbreviated by JR

Scraping related tasks will be sent to the topic abbreviated by ST

Assessment related tasks will be sent to the topic abbreviated by AT
Assessment related results will be sent to the topic abbreviated by AR

Receiving a new connection to scrape and assess:

API -> ST -> Python Service -> AT -> Golang Service -> AR -> API

Receiving a job description as a string that we need to structure:

API -> JT -> Golang Service -> JR -> API

Each service that listens to a topic has an associated "subscription".

### Dependencies

1. Docker
2. Python 3

### Setup Elasticsearch Locally
We'll use Docker to run a local instance of Elasticsearch. Run the following command in your terminal:

```
docker run -p 9200:9200 -p 9300:9300 -e "discovery.type=single-node" docker.elastic.co/elasticsearch/elasticsearch:7.14.0
```

If everything went ok, you must have Elasticsearch running on `localhost:9200`. 

Then, you can execute the `elasticsearch-indexer.ipynb` notebook. It will create the translations index and insert data from the datasets in the `/translations` folder. 

Obs.: The datasets are big so it may be slow to insert everything. If you want to test faster, you can insert only one dataset and you can stop the insertion during the operation, so it will only submit a few words.

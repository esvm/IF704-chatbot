from typing import Dict, List
from warnings import filterwarnings
from elasticsearch import Elasticsearch


filterwarnings("ignore")
es = Elasticsearch(HOST="http://localhost", PORT=9200)

ES_MAX_RESULTS = 10
ES_INDEX = "translations"

class ElasticsearchResult:
    def __init__(self, count: int, hits: List) -> None:
        self.count = count
        self.hits = hits

def buildElasticsearchQuery(sentence, language) -> Dict:
    return {
        "from": 0,
        "size": ES_MAX_RESULTS,
        "query": {
            "bool": {
                "must": [
                    {
                        "match": {
                            "language": language,
                        }
                    },
                    {
                        "match": {
                            "englishSentence": sentence
                        }
                    }
                ]
            }
        }
    }

def parseElasticsearchResult(response: Dict) -> ElasticsearchResult:
    responseHits = response["hits"]
    
    matchesCount = responseHits["total"]["value"]
    hits = responseHits["hits"]

    results = []
    for hit in hits:
        results.append(hit["_source"])

    return ElasticsearchResult(matchesCount, results)

def findTranslationsInElasticsearch(sentence, language) -> ElasticsearchResult:
    # Build ES query
    query = buildElasticsearchQuery(sentence, language)

    # Search in ES
    response = es.search(index=ES_INDEX, body=query)

    return parseElasticsearchResult(response)

def save_translation(englishSentence, translation, language):
    body = { 'englishSentence': englishSentence, 'sentence': translation, 'language': language }
    es.index(index='translations', doc_type='sentence', body=body)

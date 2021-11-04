Feature: Test Different ElasticSearch driver

    Scenario: Delete existing index
        Given no index "$DRIVER_extra_index_1" in es "extra"
        Then index "$DRIVER_extra_index_1" does not exist in es "extra"

        # Delete a non-existing index should be fine.
        Given no index "$DRIVER_extra_index_1" in es "extra"
        Then index "$DRIVER_extra_index_1" does not exist in es "extra"

    Scenario: Create index when it does not exist
        Given no index "$DRIVER_extra_index_2" in es "extra"

        When index "$DRIVER_extra_index_2" is created in es "extra"

        Then index "$DRIVER_extra_index_2" exists in es "extra"

    Scenario: Recreate index
        Given no index "$DRIVER_extra_index_3" in es "extra"
        And index "$DRIVER_extra_index_3" is created in es "extra"
        And these docs are stored in index "$DRIVER_extra_index_3" of es "extra":
        """
        [
            {
                "_id": "41",
                "_source": {
                    "handle": "item-41",
                    "name": "Item 41",
                    "locale": "en_US"
                }
            },
            {
                "_id": "42",
                "_source": {
                    "handle": "item-42",
                    "name": "Item 42",
                    "locale": "en_US"
                }
            }
        ]
        """

        When index "$DRIVER_extra_index_3" is recreated in es "extra"

        Then index "$DRIVER_extra_index_3" exists in es "extra"
        And no docs are available in index "$DRIVER_extra_index_3" of es "extra"

    Scenario: Delete and then create index
        Given no index "$DRIVER_extra_index_4" in es "extra"
        And index "$DRIVER_extra_index_4" is created in es "extra"
        And these docs are stored in index "$DRIVER_extra_index_4" of es "extra":
        """
        [
            {
                "_id": "41",
                "_source": {
                    "handle": "item-41",
                    "name": "Item 41",
                    "locale": "en_US"
                }
            },
            {
                "_id": "42",
                "_source": {
                    "handle": "item-42",
                    "name": "Item 42",
                    "locale": "en_US"
                }
            }
        ]
        """

        Given no index "$DRIVER_extra_index_4" in es "extra"
        And index "$DRIVER_extra_index_4" is created in es "extra"

        Then no docs are available in index "$DRIVER_extra_index_4" of es "extra"

    Scenario: Indexed documents are available for search
        Given index "$DRIVER_extra_index_5" is recreated in es "extra"
        And these docs are stored in index "$DRIVER_extra_index_5" of es "extra":
        """
        [
            {
                "_id": "41",
                "_source": {
                    "handle": "item-41",
                    "name": "Item 41",
                    "locale": "en_US"
                }
            },
            {
                "_id": "42",
                "_source": {
                    "handle": "item-42",
                    "name": "Item 42",
                    "locale": "en_US"
                }
            }
        ]
        """

        And only these docs are available in index "$DRIVER_extra_index_5" of es "extra":
        """
        [
            {
                "_id": "41",
                "_source": {
                    "handle": "item-41",
                    "name": "Item 41",
                    "locale": "en_US"
                },
                "_score": 1,
                "_type": "_doc"
            },
            {
                "_id": "42",
                "_source": {
                    "handle": "item-42",
                    "name": "Item 42",
                    "locale": "en_US"
                },
                "_score": 1,
                "_type": "_doc"
            }
        ]
        """

        Given no docs in index "$DRIVER_extra_index_5" of es "extra"
        Then no docs are available in index "$DRIVER_extra_index_5" of es "extra"

    Scenario: Indexed documents from a file are available for search
        Given index "$DRIVER_extra_index_6" is recreated in es "extra"
        And docs in this file are stored in index "$DRIVER_extra_index_6" of es "extra":
        """
        ../../resources/fixtures/products_en_us.json
        """

        And only docs in this file are available in index "$DRIVER_extra_index_6" of es "extra":
        """
        ../../resources/fixtures/result_en_us.json
        """

    Scenario: Search for documents by query
        Given there is index "$DRIVER_extra_index_7" in es "extra"
        And these docs are stored in index "$DRIVER_extra_index_7" of es "extra":
        """
        [
            {
                "_id": "41",
                "_source": {
                    "handle": "item-41",
                    "name": "Item 41",
                    "locale": "en_US"
                }
            },
            {
                "_id": "42",
                "_source": {
                    "handle": "item-42",
                    "name": "Item 42",
                    "locale": "en_US"
                }
            },
            {
                "_id": "43",
                "_source": {
                    "handle": "item-43",
                    "name": "Item 43",
                    "locale": "fr_FR"
                }
            }
        ]
        """

        When I search in index "$DRIVER_extra_index_7" of es "extra" with query:
        """
        {
            "query": {
                "match": {
                    "locale": "en_US"
                }
            }
        }
        """

        Then these docs are found in index "$DRIVER_extra_index_7" of es "extra":
        """
        [
            {
                "_id": "41",
                "_source": {
                    "handle": "item-41",
                    "name": "Item 41",
                    "locale": "en_US"
                },
                "_score": "<ignore-diff>",
                "_type": "_doc"
            },
            {
                "_id": "42",
                "_source": {
                    "handle": "item-42",
                    "name": "Item 42",
                    "locale": "en_US"
                },
                "_score": "<ignore-diff>",
                "_type": "_doc"
            }
        ]
        """

    Scenario: Search for documents from a file by query
        Given there is index "$DRIVER_extra_index_8" in es "extra"
        And docs in this file are stored in index "$DRIVER_extra_index_8" of es "extra":
        """
        ../../resources/fixtures/products_mixed.json
        """

        When I search in index "$DRIVER_extra_index_8" of es "extra" with query:
        """
        {
            "query": {
                "match": {
                    "locale": "en_US"
                }
            }
        }
        """

        Then docs in this file are found in index "$DRIVER_extra_index_8" of es "extra":
        """
        ../../resources/fixtures/result_en_us.json
        """

    Scenario: Create index with inline config
        Given no index "$DRIVER_extra_index_9" in es "extra"

        When index "$DRIVER_extra_index_9" is created in es "extra" with config:
        """
        {
            "mappings": {
                "properties": {
                    "age": {
                        "type": "integer"
                    },
                    "email": {
                        "type": "keyword"
                    },
                    "name": {
                        "type": "text"
                    }
                }
            }
        }
        """

        Then index "$DRIVER_extra_index_9" exists in es "extra"

    Scenario: Create index with config from a file
        Given no index "$DRIVER_extra_index_10" in es "extra"

        When index "$DRIVER_extra_index_10" is created in es "extra" with config from file:
        """
        ../../resources/fixtures/mapping.json
        """

        Then index "$DRIVER_extra_index_10" exists in es "extra"

    Scenario: Recreate index with inline config
        Given index "$DRIVER_default_index_11" is recreated in es "extra" with config:
        """
        {
            "mappings": {
                "properties": {
                    "age": {
                        "type": "integer"
                    },
                    "email": {
                        "type": "keyword"
                    },
                    "name": {
                        "type": "text"
                    }
                }
            }
        }
        """

        Then index "$DRIVER_default_index_11" exists in es "extra"

    Scenario: Recreate index with config from a file
        Given index "$DRIVER_default_index_12" is recreated in es "extra" with config from file:
        """
        ../../resources/fixtures/mapping.json
        """

        Then index "$DRIVER_default_index_12" exists in es "extra"

    Scenario: There is an index with inline config
        Given there is an index "$DRIVER_default_index_13" in es "extra" with config:
        """
        {
            "mappings": {
                "properties": {
                    "age": {
                        "type": "integer"
                    },
                    "email": {
                        "type": "keyword"
                    },
                    "name": {
                        "type": "text"
                    }
                }
            }
        }
        """

        Then index "$DRIVER_default_index_13" exists in es "extra"

    Scenario: There is an index with config from a file
        Given there is an index "$DRIVER_default_index_14" in es "extra" with config from file:
        """
        ../../resources/fixtures/mapping.json
        """

        Then index "$DRIVER_default_index_14" exists in es "extra"

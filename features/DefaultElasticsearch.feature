Feature: Test Different ElasticSearch driver

    Scenario: Delete existing index
        Given no index "$DRIVER_default_index_1"
        Then index "$DRIVER_default_index_1" does not exist

        # Delete a non-existing index should be fine.
        Given no index "$DRIVER_default_index_1"
        Then index "$DRIVER_default_index_1" does not exist

    Scenario: Create index when it does not exist
        Given no index "$DRIVER_default_index_2"

        When index "$DRIVER_default_index_2" is created

        Then index "$DRIVER_default_index_2" exists

    Scenario: Recreate index
        Given no index "$DRIVER_default_index_3"
        And index "$DRIVER_default_index_3" is created
        And these docs are stored in index "$DRIVER_default_index_3":
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

        When index "$DRIVER_default_index_3" is recreated

        Then index "$DRIVER_default_index_3" exists
        And no docs are available in index "$DRIVER_default_index_3"

    Scenario: Delete and then create index
        Given no index "$DRIVER_default_index_4"
        And index "$DRIVER_default_index_4" is created
        And these docs are stored in index "$DRIVER_default_index_4":
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

        Given no index "$DRIVER_default_index_4"
        And index "$DRIVER_default_index_4" is created

        Then no docs are available in index "$DRIVER_default_index_4"

    Scenario: Indexed documents are available for search
        Given index "$DRIVER_default_index_5" is recreated
        And these docs are stored in index "$DRIVER_default_index_5":
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

        And only these docs are available in index "$DRIVER_default_index_5":
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

        Given no docs in index "$DRIVER_default_index_5"
        Then no docs are available in index "$DRIVER_default_index_5"

    Scenario: Indexed documents from a file are available for search
        Given index "$DRIVER_default_index_6" is recreated
        And docs in this file are stored in index "$DRIVER_default_index_6":
        """
        ../../resources/fixtures/products_en_us.json
        """

        And only docs in this file are available in index "$DRIVER_default_index_6":
        """
        ../../resources/fixtures/result_en_us.json
        """

    Scenario: Search for documents by query
        Given there is index "$DRIVER_default_index_7"
        And these docs are stored in index "$DRIVER_default_index_7":
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

        When I search in index "$DRIVER_default_index_7" with query:
        """
        {
            "query": {
                "match": {
                    "locale": "en_US"
                }
            }
        }
        """

        Then these docs are found in index "$DRIVER_default_index_7":
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
        Given there is index "$DRIVER_default_index_8"
        And docs in this file are stored in index "$DRIVER_default_index_8":
        """
        ../../resources/fixtures/products_mixed.json
        """

        When I search in index "$DRIVER_default_index_8" with query:
        """
        {
            "query": {
                "match": {
                    "locale": "en_US"
                }
            }
        }
        """

        Then docs in this file are found in index "$DRIVER_default_index_8":
        """
        ../../resources/fixtures/result_en_us.json
        """

    Scenario: Create index with inline config
        Given no index "$DRIVER_default_index_9"

        When index "$DRIVER_default_index_9" is created with config:
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

        Then index "$DRIVER_default_index_9" exists

    Scenario: Create index with config from a file
        Given no index "$DRIVER_default_index_10"

        When index "$DRIVER_default_index_10" is created with config from file:
        """
        ../../resources/fixtures/mapping.json
        """

        Then index "$DRIVER_default_index_10" exists

    Scenario: Recreate index with inline config
        Given index "$DRIVER_default_index_11" is recreated with config:
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

        Then index "$DRIVER_default_index_11" exists

    Scenario: Recreate index with config from a file
        Given index "$DRIVER_default_index_12" is recreated with config from file:
        """
        ../../resources/fixtures/mapping.json
        """

        Then index "$DRIVER_default_index_12" exists

    Scenario: There is an index with inline config
        Given there is an index "$DRIVER_default_index_13" with config:
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

        Then index "$DRIVER_default_index_13" exists

    Scenario: There is an index with config from a file
        Given there is an index "$DRIVER_default_index_14" with config from file:
        """
        ../../resources/fixtures/mapping.json
        """

        Then index "$DRIVER_default_index_14" exists

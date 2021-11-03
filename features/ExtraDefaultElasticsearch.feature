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

    Scenario: Search for documents by query
        Given there is index "$DRIVER_extra_index_6" in es "extra"
        And these docs are stored in index "$DRIVER_extra_index_6" of es "extra":
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

        When I search in index "$DRIVER_extra_index_6" of es "extra" with query:
        """
        {
            "query": {
                "match": {
                    "locale": "en_US"
                }
            }
        }
        """

        Then these docs are found in index "$DRIVER_extra_index_6" of es "extra":
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

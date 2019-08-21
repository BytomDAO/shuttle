import requests
import json

def get_testnet_coins(receiver_str):
    transaction_dict = {
        "actions":[
            {
                "account_id":"0KN9JNBA00A02",
                "amount":1020000000,
                "asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
                "type":"spend_account",
                "use_unconfirmed":True
            },
            {
                "amount":1000000000,
                "asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
                "address":receiver_str,
                "type":"control_address"
            }
        ],
        "ttl":0,
        "time_range":0
    }
    transaction_json = json.dumps(transaction_dict)
    headers = {
        "content-type": "application/json",
        "accept": "application/json"
    }
    build_url = "http://127.0.0.1:9888/build-transaction"
    response = requests.post(build_url, data=transaction_json).json()
    built_transaction_dict = {
        "password": "12345",
        "transaction": response['data']
    }
    built_transaction_json = json.dumps(built_transaction_dict)
    sign_url = "http://127.0.0.1:9888/sign-transaction"
    response = requests.post(sign_url, headers=headers, data=built_transaction_json).json()
    signed_transaction_dict = {
        "raw_transaction": response['data']['transaction']['raw_transaction']
    }
    signed_transaction_json = json.dumps(signed_transaction_dict)
    submit_url = "http://127.0.0.1:9888/submit-transaction"
    response = requests.post(submit_url, headers=headers, data=signed_transaction_json).json()
    return {
        "tx_id": response['data']['tx_id']
    }

def get_gm_testnet_coins(receiver_str):
    transaction_dict = {
        "actions":[
            {
                "account_id":"0KN8NES7G0A02",
                "amount":1020000000,
                "asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
                "type":"spend_account",
                "use_unconfirmed":True
            },
            {
                "amount":1000000000,
                "asset_id":"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
                "address":receiver_str,
                "type":"control_address"
            }
        ],
        "ttl":0,
        "time_range":0
    }
    transaction_json = json.dumps(transaction_dict)
    headers = {
        "content-type": "application/json",
        "accept": "application/json"
    }
    build_url = "http://127.0.0.1:9889/build-transaction"
    response = requests.post(build_url, data=transaction_json).json()
    built_transaction_dict = {
        "password": "12345",
        "transaction": response['data']
    }
    built_transaction_json = json.dumps(built_transaction_dict)
    sign_url = "http://127.0.0.1:9889/sign-transaction"
    response = requests.post(sign_url, headers=headers, data=built_transaction_json).json()
    signed_transaction_dict = {
        "raw_transaction": response['data']['transaction']['raw_transaction']
    }
    signed_transaction_json = json.dumps(signed_transaction_dict)
    submit_url = "http://127.0.0.1:9889/submit-transaction"
    response = requests.post(submit_url, headers=headers, data=signed_transaction_json).json()
    return {
        "tx_id": response['data']['tx_id']
    }

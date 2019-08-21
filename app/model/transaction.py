import requests
import json
from _pysha3 import sha3_256
from app.model import receiver

# submit_transaction broadcast raw transaction
# raw_transaction_str is signed transaction,
# network_str is mainnet or testnet
# test data 1:
#   raw_transaction_str: 070100010160015e0873eddd68c4ba07c9410984799928288ae771bdccc6d974e72c95727813461fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8094ebdc030101160014052620b86a6d5e07311d5019dffa3864ccc8a6bd630240312a052f36efb9826aa1021ec91bc6f125dd07f9c4bff87014612069527e15246518806b654d57fff8b6fe91866a19d5a2fb63a5894335fce92a7b4a7fcd340720e87ca3acdebdcad9a1d0f2caecf8ce0dbfc73d060807a210c6f225488347961402013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8082eee0020116001418028ef4f8b8c278907864a1977a5ee6707b2a6b00013cffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80b8b872011600142935e4869d0317d9701c80a02ecf888143cb9dd200
#   network_str: testnet
def submit_transaction(raw_transaction_str, network_str):
    raw_transaction_dict = {
        "transaction": raw_transaction_str
    }
    raw_transaction_json = json.dumps(raw_transaction_dict)
    headers = {
        "content-type": "application/json",
        "accept": "application/json"
    }
    if network_str == "mainnet":
        url = "https://blockmeta.com/api/v2/broadcast-transaction"
    else:
        url = "https://blockmeta.com/api/wisdom/broadcast-transaction"
    response = requests.post(url, headers=headers, data=raw_transaction_json)
    return {
        "response": response.text[:-1]
    }


def decode_raw_transaction(raw_transaction_str):
    raw_transaction_dict = {
        "raw_transaction": raw_transaction_str
    }
    raw_transaction_json = json.dumps(raw_transaction_dict)
    headers = {
        "content-type": "application/json",
        "accept": "application/json"
    }
    url = 'http://127.0.0.1:9888/decode-raw-transaction'
    response = requests.post(url, headers=headers, data=raw_transaction_json)
    return {
        "response": response.text[:-1]
    }


def get_uvarint(uvarint_str):
    uvarint_bytes = bytes.fromhex(uvarint_str)
    x, s, i = 0, 0, 0
    while True:
        b = uvarint_bytes[i]
        if b < 0x80:
            if i > 9 or i == 9 and b > 1:
                return "overflow"
            return x | int(b) << s, i + 1
        x |= int(b & 0x7f) << s
        s += 7
        i += 1


'''
get_spend_output_id create tx_input spend output id
test data 1:
  source_id_hexstr: 28b7b53d8dc90006bf97e0a4eaae2a72ec3d869873188698b694beaf20789f21
  asset_id_hexstr: ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
  amount_int: 41250000000
  source_position_int: 0
  vmversion_int: 1
  control_program_hexstr: 00149335b1cbd4a77b78e33315a0ed10a95b12e7ca48
  spend_output_id_hexstr: f229ec6f403d586dc87aa2546bbe64c5f7b5f46eb13c6ee4823d03bc88a7cf17
'''
def get_spend_output_id(source_id_hexstr, asset_id_hexstr, amount_int, source_position_int, vmversion_int, control_program_hexstr):
    amount_hexstr = amount_int.to_bytes(8, byteorder='little').hex()
    source_position_hexstr = source_position_int.to_bytes(8, byteorder='little').hex()
    vmversion_hexstr = vmversion_int.to_bytes(8, byteorder='little').hex()
    cp_length_int = len(control_program_hexstr) // 2
    cp_length_hexstr = cp_length_int.to_bytes((cp_length_int.bit_length() + 7) // 8, byteorder='little').hex()
    sc_hexstr = source_id_hexstr + asset_id_hexstr + amount_hexstr + source_position_hexstr + vmversion_hexstr + cp_length_hexstr +  control_program_hexstr
    innerhash_bytes = sha3_256(bytes.fromhex(sc_hexstr)).digest()
    spend_bytes = b'entryid:output1:' + innerhash_bytes
    spend_output_id_hexstr = sha3_256(spend_bytes).hexdigest()
    return spend_output_id_hexstr

'''
get_input_id create tx input_id
test data 1:
    spend_output_id_hexstr: f229ec6f403d586dc87aa2546bbe64c5f7b5f46eb13c6ee4823d03bc88a7cf17
    input_id_hexstr: 6e3f378ed844b143a335e306f4ba26746157589c87e8fc8cba6463c566c56768
'''
def get_input_id(spend_output_id_hexstr):
    innerhash_bytes = sha3_256(bytes.fromhex(spend_output_id_hexstr)).digest()
    input_id_hexstr = sha3_256(b'entryid:spend1:' + innerhash_bytes).hexdigest()
    return input_id_hexstr


def get_mux_id(prepare_mux_hexstr):
    innerhash_bytes = sha3_256(bytes.fromhex(prepare_mux_hexstr)).digest()
    mux_id_hexstr = sha3_256(b'entryid:mux1:' + innerhash_bytes).hexdigest()
    return mux_id_hexstr


def get_output_id(prepare_output_id_hexstr):
    innerhash_bytes = sha3_256(bytes.fromhex(prepare_output_id_hexstr)).digest()
    output_id_hexstr = sha3_256(b'entryid:output1:' + innerhash_bytes).hexdigest()
    return output_id_hexstr


def get_tx_id(prepare_tx_id_hexstr):
    innerhash_bytes = sha3_256(bytes.fromhex(prepare_tx_id_hexstr)).digest()
    tx_id_hexstr = sha3_256(b'entryid:txheader:' + innerhash_bytes).hexdigest()
    return tx_id_hexstr


def get_issue_input_id(prepare_issue_hexstr):
    innerhash_bytes = sha3_256(bytes.fromhex(prepare_issue_hexstr)).digest()
    tx_id_hexstr = sha3_256(b'entryid:issuance1:' + innerhash_bytes).hexdigest()
    return tx_id_hexstr


def get_coinbase_input_id(prepare_coinbase_input_id_hexstr):
    innerhash_bytes = sha3_256(bytes.fromhex(prepare_coinbase_input_id_hexstr)).digest()
    coinbase_input_id_hexstr = sha3_256(b'entryid:coinbase1:' + innerhash_bytes).hexdigest()
    return coinbase_input_id_hexstr


'''
decode_raw_tx decode raw transaction
testdata 1:
    raw_transaction_str: 070100010161015f28b7b53d8dc90006bf97e0a4eaae2a72ec3d869873188698b694beaf20789f21ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8099c4d5990100011600149335b1cbd4a77b78e33315a0ed10a95b12e7ca48630240897e2d9d24a3b5faaed0579dee7597b401491595675f897504f8945b29d836235bd2fca72a3ad0cae814628973ebcd142d9d6cc92d0b2571b69e5370a98a340c208cb7fb3086f58db9a31401b99e8c658be66134fb9034de1d5c462679270b090702013effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80f9f8bc98010116001406ce4b689ba026ffd3a7ca65d1d059546d4b78a000013dffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff80c6868f01011600147929ef91997c827bebf60fa608f876ea27523c4700
    network_str: solonet
    transaction: 
{
  "fee": 20000000,
  "inputs": [
    {
      "address": "sm1qjv6mrj755aah3cenzksw6y9ftvfw0jjgk0l2mw",
      "amount": 41250000000,
      "asset_definition": {},
      "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
      "control_program": "00149335b1cbd4a77b78e33315a0ed10a95b12e7ca48",
      "input_id": "6e3f378ed844b143a335e306f4ba26746157589c87e8fc8cba6463c566c56768",
      "spent_output_id": "f229ec6f403d586dc87aa2546bbe64c5f7b5f46eb13c6ee4823d03bc88a7cf17",
      "type": "spend",
      "witness_arguments": [
        "897e2d9d24a3b5faaed0579dee7597b401491595675f897504f8945b29d836235bd2fca72a3ad0cae814628973ebcd142d9d6cc92d0b2571b69e5370a98a340c",
        "8cb7fb3086f58db9a31401b99e8c658be66134fb9034de1d5c462679270b0907"
      ]
    }
  ],
  "outputs": [
    {
      "address": "sm1qqm8yk6ym5qn0l5a8efjar5ze23k5k79qnvtslj",
      "amount": 40930000000,
      "asset_definition": {},
      "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
      "control_program": "001406ce4b689ba026ffd3a7ca65d1d059546d4b78a0",
      "id": "74c73266730d3c6ea32e8667ef9b867068736b84be240fe9fef205fa68bb7b95",
      "position": 0,
      "type": "control"
    },
    {
      "address": "sm1q0y57lyve0jp8h6lkp7nq37rkagn4y0z8hvh6kq",
      "amount": 300000000,
      "asset_definition": {},
      "asset_id": "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
      "control_program": "00147929ef91997c827bebf60fa608f876ea27523c47",
      "id": "f115a833d0c302a5006032858a7ed3987f0feb2daf2a9f849384950e4766af51",
      "position": 1,
      "type": "control"
    }
  ],
  "size": 333,
  "time_range": 0,
  "tx_id": "814a73dd57bae67c604f9cbc696cbc42035577423408cb9267136ed971e2bf63",
  "version": 1
}
'''
def decode_raw_tx(raw_transaction_str, network_str):
    tx = {
        "fee": 0,
        "inputs": [],
        "outputs": [],
        "size": 0,
        "time_range": 0,
        "tx_id": "",
        "version": 0
    }
    tx['fee'] = 0
    tx['size'] = len(raw_transaction_str) // 2
    length = 0
    offset = 2
    tx['version'], length = get_uvarint(raw_transaction_str[offset:offset+18])
    offset = offset + 2 * length
    tx['time_range'], length = get_uvarint(raw_transaction_str[offset:offset+18])
    offset = offset + 2 * length
    tx_input_amount, length = get_uvarint(raw_transaction_str[offset:offset+8])
    offset = offset + 2 * length
    prepare_mux_hexstr = (tx_input_amount).to_bytes((tx_input_amount.bit_length() + 7) // 8, 'little').hex()
    prepare_tx_id_hexstr = (tx['version']).to_bytes(8, 'little').hex() + (tx['time_range']).to_bytes(8, 'little').hex()
    for _ in range(tx_input_amount):
        _, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        _, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        input_type = int(raw_transaction_str[offset:offset+2], 16)
        offset += 2
        if input_type == 0: # issue
            tx_input = {
                "amount": 0,
                "asset_definition": "", # TODO:fix it!!
                "asset_id": "",
                "input_id": "",
                "issuance_program": "",
                "type": "",
                "witness_arguments": []
            }
            tx_input['type'] = "issue"
            _, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            nonce = raw_transaction_str[offset:offset+16]
            offset += 16
            nonce_hash_hexstr = sha3_256(bytes.fromhex(nonce)).hexdigest()
            tx_input['asset_id'] = raw_transaction_str[offset:offset+64]
            offset += 64
            tx_input['amount'], length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            _, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            asset_definition_size, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            tx_input['asset_definition'] = bytes.fromhex(raw_transaction_str[offset:offset+2*asset_definition_size]).decode()
            offset = offset + 2 * asset_definition_size
            _, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            issuance_program_length, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            tx_input['issuance_program'] = raw_transaction_str[offset:offset+2*issuance_program_length]
            offset = offset + 2 * issuance_program_length
            witness_arguments_amount, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            for _ in range(witness_arguments_amount):
                argument_length, length = get_uvarint(raw_transaction_str[offset:offset+18])
                offset = offset + 2 * length
                argument = raw_transaction_str[offset:offset+2*argument_length]
                offset = offset + 2 * argument_length
                tx_input['witness_arguments'].append(argument)
            prepare_issue_hexstr = nonce_hash_hexstr + tx_input['asset_id'] + (tx_input['amount']).to_bytes(8, byteorder='little').hex()
            tx_input['input_id'] = get_issue_input_id(prepare_issue_hexstr)
            tx['inputs'].append(tx_input)
            prepare_mux_hexstr += tx_input['input_id'] + tx_input['asset_id'] + (tx_input['amount']).to_bytes(8, byteorder='little').hex() + '0000000000000000'
            prepare_mux_hexstr += '0100000000000000' + '0151'
            mux_id_hexstr = get_mux_id(prepare_mux_hexstr)
        elif input_type == 1: # spend
            tx_input = {
                "address": "",
                "amount": 0,
                "asset_definition": {},
                "asset_id": "",
                "control_program": "",
                "input_id": "",
                "spent_output_id": "",
                "type": "",
                "witness_arguments": []
            }
            tx_input['type'] = "spend"
            _, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            source_id = raw_transaction_str[offset:offset+64]
            offset += 64
            tx_input['asset_id'] = raw_transaction_str[offset:offset+64]
            offset += 64
            tx_input['amount'], length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            tx['fee'] += tx_input['amount']
            source_positon, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            vmversion, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            control_program_length, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            tx_input['control_program'] = raw_transaction_str[offset:offset+2*control_program_length]
            offset = offset + 2 * control_program_length
            tx_input['address'] = receiver.create_address(tx_input['control_program'], network_str)['address']
            _, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            witness_arguments_amount, length = get_uvarint(raw_transaction_str[offset:offset+18])
            offset = offset + 2 * length
            if witness_arguments_amount == 1:
                offset = offset + 2
                tx_input['witness_arguments'] = None
            else: 
                for _ in range(witness_arguments_amount):
                    argument_length, length = get_uvarint(raw_transaction_str[offset:offset+18])
                    offset = offset + 2 * length
                    argument = raw_transaction_str[offset:offset+2*argument_length]
                    offset = offset + 2 * argument_length
                    tx_input['witness_arguments'].append(argument)
            tx_input['spent_output_id'] = get_spend_output_id(source_id, tx_input['asset_id'], tx_input['amount'], source_positon, vmversion, tx_input['control_program'])
            tx_input['input_id'] = get_input_id(tx_input['spent_output_id'])
            tx['inputs'].append(tx_input)
            prepare_mux_hexstr += tx_input['input_id'] + tx_input['asset_id'] + (tx_input['amount']).to_bytes(8, byteorder='little').hex() + '0000000000000000'
            prepare_mux_hexstr += '0100000000000000' + '0151'
            mux_id_hexstr = get_mux_id(prepare_mux_hexstr)
        elif input_type == 2: # coinbase
            tx_input = {
                "amount": 0,
                "arbitrary": "",
                "asset_definition": {},
                "asset_id": "0000000000000000000000000000000000000000000000000000000000000000",
                "input_id": "",
                "type": "",
                "witness_arguments": []
            }
            tx_input['type'] = "coinbase"
            arbitrary_length, length = get_uvarint(raw_transaction_str[offset:offset+18])
            prepare_coinbase_input_id_hexstr = raw_transaction_str[offset:offset+2*length]
            offset = offset + 2 * length
            tx_input['arbitrary'] = raw_transaction_str[offset:offset+2*arbitrary_length]
            prepare_coinbase_input_id_hexstr += tx_input['arbitrary']
            offset = offset + 2 * arbitrary_length
            tx_input['input_id'] = get_coinbase_input_id(prepare_coinbase_input_id_hexstr)
            offset = offset + 2
            tx['inputs'].append(tx_input)
            prepare_mux_hexstr += tx_input['input_id'] + 'ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff'
    tx_output_amount, length = get_uvarint(raw_transaction_str[offset:offset+18])
    offset = offset + 2 * length
    prepare_tx_id_hexstr += (tx_output_amount).to_bytes((tx_output_amount.bit_length() + 7) // 8, 'little').hex()
    for i in range(tx_output_amount):
        tx_output = {
            "address": "",
            "amount": 0,
            "asset_definition": {},
            "asset_id": "",
            "control_program": "",
            "id": "",
            "position": 0,
            "type": ""
        }
        tx_output['position'] = i
        _, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        _, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        tx_output['asset_id'] = raw_transaction_str[offset:offset+64]
        offset = offset + 64
        tx_output['amount'], length = get_uvarint(raw_transaction_str[offset:offset+18])
        if tx_input['type'] == "coinbase":
            prepare_mux_hexstr = prepare_mux_hexstr + (tx_output['amount']).to_bytes(8, byteorder='little').hex() + '0000000000000000' + '0100000000000000' + '0151'
            mux_id_hexstr = get_mux_id(prepare_mux_hexstr)
        offset = offset + 2 * length
        if tx_output['asset_id'] == 'ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff':
            tx['fee'] -= tx_output['amount']
        _, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        control_program_length, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        tx_output['control_program'] = raw_transaction_str[offset:offset+2*control_program_length]
        offset = offset + 2 * control_program_length
        tx_output['address'] = receiver.create_address(tx_output['control_program'], network_str)['address']
        _, length = get_uvarint(raw_transaction_str[offset:offset+18])
        offset = offset + 2 * length
        prepare_output_id_hexstr = mux_id_hexstr + tx_output['asset_id'] + (tx_output['amount']).to_bytes(8, byteorder='little').hex() + (i).to_bytes(8, byteorder='little').hex() + '0100000000000000' + (control_program_length).to_bytes((control_program_length.bit_length() + 7) // 8, 'little').hex() + tx_output['control_program']
        tx_output['id'] = get_output_id(prepare_output_id_hexstr)
        prepare_tx_id_hexstr += tx_output['id']
        tx_output['type'] = 'control'
        tx['outputs'].append(tx_output)
    if tx_input['type'] == "coinbase":
        tx['fee'] = 0
    tx['tx_id'] = get_tx_id(prepare_tx_id_hexstr)
    return tx

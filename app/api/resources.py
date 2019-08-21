from flask_restful import inputs
from flask_restful import Resource
from flask_restful import reqparse
from app.model.hello import get_hello_result
from app.model.signature import sign
from app.model.signature import verify
from app.model.key import create_entropy
from app.model.key import entropy_to_mnemonic
from app.model.key import mnemonic_to_seed
from app.model.key import seed_to_root_xprv
from app.model.key import xprv_to_expanded_private_key
from app.model.key import xpub_to_public_key
from app.model.key import xpub_verify
from app.model.key import xprv_sign
from app.model.key import xprv_to_xpub
from app.model.key import xprv_to_child_xprv
from app.model.key import xpub_to_child_xpub
from app.model.receiver import create_P2WPKH_program
from app.model.receiver import create_address
from app.model.receiver import get_path_from_index
from app.model.receiver import create_qrcode_base64
from app.model.transaction import submit_transaction
from app.model.key import create_new_key
from app.model.receiver import create_new_address
from app.model.transaction import decode_raw_transaction
from app.model.key_gm import get_gm_root_xprv
from app.model.key_gm import get_gm_xpub
from app.model.key_gm import get_gm_xprv
from app.model.key_gm import get_gm_public_key
from app.model.key_gm import get_gm_child_xprv
from app.model.key_gm import get_gm_child_xpub
from app.model.key_gm import gm_xpub_verify
from app.model.key_gm import gm_xprv_sign
from app.model.receiver import get_gm_P2WPKH_program
from app.model.receiver import get_gm_address
from app.model.key_gm import get_gm_new_key
from app.model.receiver import get_gm_new_address
from app.model.transaction import decode_raw_tx
from app.model.faucet import get_testnet_coins
from app.model.faucet import get_gm_testnet_coins


parser = reqparse.RequestParser()
parser.add_argument('private_key_str', type=str)
parser.add_argument('message_str', type=str)
parser.add_argument('public_key_str', type=str)
parser.add_argument('signature_str', type=str)
parser.add_argument('entropy_str', type=str)
parser.add_argument('mnemonic_str', type=str)
parser.add_argument('seed_str', type=str)
parser.add_argument('xprv_str', type=str)
parser.add_argument('xpub_str', type=str)
parser.add_argument('path_list', type=str, action='append')
parser.add_argument('account_index_int', type=int)
parser.add_argument('address_index_int', type=int)
parser.add_argument('change_bool', type=inputs.boolean)
parser.add_argument('control_program_str', type=str)
parser.add_argument('network_str', type=str)
parser.add_argument('raw_transaction_str', type=str)
parser.add_argument('address_str', type=str)
parser.add_argument('s', type=str)
parser.add_argument('receiver_str', type=str)

class Hello(Resource):

    def get(self, content):
        return get_hello_result(content)

class Sign(Resource):

    def post(self):
        args = parser.parse_args()
        private_key = args.get('private_key_str')
        message = args.get('message_str')
        sig = sign(private_key, message)
        return sig

class Verify(Resource):
    
    def post(self):
        args = parser.parse_args()
        public_key = args.get('public_key_str')
        signature = args.get('signature_str')
        message = args.get('message_str')
        result = verify(public_key, signature, message)
        return result

class Create_Entropy(Resource):

    def post(self):
        entropy = create_entropy()
        return entropy

class Entropy_To_Mnemonic(Resource):

    def post(self):
        args = parser.parse_args()
        entropy_str = args.get('entropy_str')
        mnemonic_str = entropy_to_mnemonic(entropy_str)
        return mnemonic_str

class Mnemonic_To_Seed(Resource):

    def post(self):
        args = parser.parse_args()
        mnemonic_str = args.get('mnemonic_str')
        seed_str = mnemonic_to_seed(mnemonic_str)
        return seed_str

class Seed_To_Root_Xprv(Resource):

    def post(self):
        args = parser.parse_args()
        seed_str = args.get('seed_str')
        root_xprv_str = seed_to_root_xprv(seed_str)
        return root_xprv_str

class Xprv_To_Expanded_Private_Key(Resource):

    def post(self):
        args = parser.parse_args()
        xprv_str = args.get('xprv_str')
        expanded_private_key_str = xprv_to_expanded_private_key(xprv_str)
        return expanded_private_key_str

class Xpub_To_Public_Key(Resource):

    def post(self):
        args = parser.parse_args()
        xpub_str = args.get('xpub_str')
        public_key_str = xpub_to_public_key(xpub_str)
        return public_key_str

class Xpub_Verify(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        message = args.get('message_str')
        signature = args.get('signature_str')
        result = xpub_verify(xpub, message, signature)
        return result

class Xprv_Sign(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        message = args.get('message_str')
        signature = xprv_sign(xprv, message)
        return signature

class Xprv_To_Xpub(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        xpub = xprv_to_xpub(xprv)
        return xpub

class Xprv_To_Child_Xprv(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        path = args.get('path_list')
        child_xprv = xprv_to_child_xprv(xprv, path)
        return child_xprv

class Xpub_To_Child_Xpub(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        path = args.get('path_list')
        child_xpub = xpub_to_child_xpub(xpub, path)
        return child_xpub

class Get_Path_From_Index(Resource):

    def post(self):
        args = parser.parse_args()
        account_index = args.get('account_index_int')
        address_index = args.get('address_index_int')
        change = args.get('change_bool')
        path_list = get_path_from_index(account_index, address_index, change)
        return path_list


class Create_P2WPKH_Program(Resource):

    def post(self):
        args = parser.parse_args()
        account_index = args.get('account_index_int')
        address_index = args.get('address_index_int')
        change = args.get('change_bool')
        xpub = args.get('xpub_str')
        control_program = create_P2WPKH_program(account_index, address_index, change, xpub)
        return control_program

class Create_Address(Resource):

    def post(self):
        args = parser.parse_args()
        control_program = args.get('control_program_str')
        network = args.get('network_str')
        address = create_address(control_program, network)
        return address

class Submit_Transaction(Resource):

    def post(self):
        args = parser.parse_args()
        raw_transaction = args.get('raw_transaction_str')
        network = args.get('network_str')
        response = submit_transaction(raw_transaction, network)
        return response

class Create_QRcode_Base64(Resource):

    def post(self):
        args = parser.parse_args()
        s = args.get('s')
        img = create_qrcode_base64(s)
        return img

class Create_New_Key(Resource):

    def post(self):
        args = parser.parse_args()
        entropy = args.get('entropy_str')
        mnemonic = args.get('mnemonic_str')
        response = create_new_key(entropy, mnemonic)
        return response


class Create_New_Address(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        account_index = args.get('account_index_int')
        address_index = args.get('address_index_int')
        change = args.get('change_bool')
        network = args.get('network_str')
        response = create_new_address(xpub, account_index, address_index, change, network)
        return response

class Decode_Raw_Transaction(Resource):

    def post(self):
        args = parser.parse_args()
        raw_transaction = args.get('raw_transaction_str')
        response = decode_raw_transaction(raw_transaction)
        return response

class Get_Gm_Root_Xprv(Resource):

    def post(self):
        args = parser.parse_args()
        seed = args.get('seed_str')
        root_xprv = get_gm_root_xprv(seed)
        return root_xprv

class Get_Gm_Xpub(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        xpub = get_gm_xpub(xprv)
        return xpub

class Get_Gm_Xprv(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        xprv = get_gm_xprv(xprv)
        return xprv

class Get_Gm_Public_Key(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        public_key = get_gm_public_key(xpub)
        return public_key

class Get_Gm_Child_Xprv(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        path = args.get('path_list')
        child_xprv = get_gm_child_xprv(xprv, path)
        return child_xprv

class Get_Gm_Child_Xpub(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        path = args.get('path_list')
        child_xpub = get_gm_child_xpub(xpub, path)
        return child_xpub

class Gm_Xpub_Verify(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        message = args.get('message_str')
        signature = args.get('signature_str')
        result = gm_xpub_verify(xpub, message, signature)
        return result

class Gm_Xprv_Sign(Resource):

    def post(self):
        args = parser.parse_args()
        xprv = args.get('xprv_str')
        message = args.get('message_str')
        signature = gm_xprv_sign(xprv, message)
        return signature

class Get_Gm_P2WPKH_Program(Resource):

    def post(self):
        args = parser.parse_args()
        account_index = args.get('account_index_int')
        address_index = args.get('address_index_int')
        change = args.get('change_bool')
        xpub = args.get('xpub_str')
        control_program = get_gm_P2WPKH_program(account_index, address_index, change, xpub)
        return control_program

class Get_Gm_Address(Resource):

    def post(self):
        args = parser.parse_args()
        control_program = args.get('control_program_str')
        network = args.get('network_str')
        address = get_gm_address(control_program, network)
        return address

class Get_Gm_New_Key(Resource):

    def post(self):
        args = parser.parse_args()
        entropy = args.get('entropy_str')
        mnemonic = args.get('mnemonic_str')
        response = get_gm_new_key(entropy, mnemonic)
        return response

class Get_Gm_New_Address(Resource):

    def post(self):
        args = parser.parse_args()
        xpub = args.get('xpub_str')
        account_index = args.get('account_index_int')
        address_index = args.get('address_index_int')
        change = args.get('change_bool')
        network = args.get('network_str')
        response = get_gm_new_address(xpub, account_index, address_index, change, network)
        return response

class Decode_Raw_Tx(Resource):

    def post(self):
        args = parser.parse_args()
        network = args.get('network_str')
        raw_transaction = args.get('raw_transaction_str')
        tx = decode_raw_tx(raw_transaction, network)
        return tx

class Get_Testnet_Coins(Resource):

    def post(self):
        args = parser.parse_args()
        receiver = args.get('receiver_str')
        tx_id = get_testnet_coins(receiver)
        return tx_id

class Get_Gm_Testnet_Coins(Resource):

    def post(self):
        args = parser.parse_args()
        receiver = args.get('receiver_str')
        tx_id = get_gm_testnet_coins(receiver)
        return tx_id
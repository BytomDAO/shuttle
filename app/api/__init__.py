from flask import Blueprint
from flask_restful import Api

from app.api.resources import Hello
from app.api.resources import Sign
from app.api.resources import Verify
from app.api.resources import Create_Entropy
from app.api.resources import Entropy_To_Mnemonic
from app.api.resources import Mnemonic_To_Seed
from app.api.resources import Seed_To_Root_Xprv
from app.api.resources import Xprv_To_Expanded_Private_Key
from app.api.resources import Xpub_To_Public_Key
from app.api.resources import Xpub_Verify
from app.api.resources import Xprv_Sign
from app.api.resources import Xprv_To_Xpub
from app.api.resources import Xprv_To_Child_Xprv
from app.api.resources import Xpub_To_Child_Xpub
from app.api.resources import Create_P2WPKH_Program
from app.api.resources import Create_Address
from app.api.resources import Get_Path_From_Index
from app.api.resources import Submit_Transaction
from app.api.resources import Create_QRcode_Base64
from app.api.resources import Create_New_Key
from app.api.resources import Create_New_Address
from app.api.resources import Decode_Raw_Transaction
from app.api.resources import Get_Gm_Root_Xprv
from app.api.resources import Get_Gm_Xpub
from app.api.resources import Get_Gm_Xprv
from app.api.resources import Get_Gm_Public_Key
from app.api.resources import Get_Gm_Child_Xprv
from app.api.resources import Get_Gm_Child_Xpub
from app.api.resources import Gm_Xpub_Verify
from app.api.resources import Gm_Xprv_Sign
from app.api.resources import Get_Gm_P2WPKH_Program
from app.api.resources import Get_Gm_Address
from app.api.resources import Get_Gm_New_Key
from app.api.resources import Get_Gm_New_Address
from app.api.resources import Decode_Raw_Tx
from app.api.resources import Get_Testnet_Coins
from app.api.resources import Get_Gm_Testnet_Coins


blueprint = Blueprint('api', __name__, url_prefix='/api/v1')
api = Api(blueprint)

api.add_resource(Hello, '/hello/<string:content>')
api.add_resource(Sign, '/sign')
api.add_resource(Verify, '/verify')
api.add_resource(Create_Entropy, '/create_entropy')
api.add_resource(Entropy_To_Mnemonic, '/entropy_to_mnemonic')
api.add_resource(Mnemonic_To_Seed, '/mnemonic_to_seed')
api.add_resource(Seed_To_Root_Xprv, '/seed_to_root_xprv')
api.add_resource(Xprv_To_Expanded_Private_Key, '/xprv_to_expanded_private_key')
api.add_resource(Xpub_To_Public_Key, '/xpub_to_public_key')
api.add_resource(Xpub_Verify, '/xpub_verify')
api.add_resource(Xprv_Sign, '/xprv_sign')
api.add_resource(Xprv_To_Xpub, '/xprv_to_xpub')
api.add_resource(Xprv_To_Child_Xprv, '/xprv_to_child_xprv')
api.add_resource(Xpub_To_Child_Xpub, '/xpub_to_child_xpub')
api.add_resource(Create_P2WPKH_Program, '/create_P2WPKH_program')
api.add_resource(Create_Address, '/create_address')
api.add_resource(Get_Path_From_Index, '/get_path_from_index')
api.add_resource(Submit_Transaction, '/submit_transaction')
api.add_resource(Create_QRcode_Base64, '/create_qrcode_base64')
api.add_resource(Create_New_Key, '/create_new_key')
api.add_resource(Create_New_Address, '/create_new_address')
api.add_resource(Decode_Raw_Transaction, '/decode_raw_transaction')
api.add_resource(Get_Gm_Root_Xprv, '/get_gm_root_xprv')
api.add_resource(Get_Gm_Xpub, '/get_gm_xpub')
api.add_resource(Get_Gm_Xprv, '/get_gm_xprv')
api.add_resource(Get_Gm_Public_Key, '/get_gm_public_key')
api.add_resource(Get_Gm_Child_Xprv, '/get_gm_child_xprv')
api.add_resource(Get_Gm_Child_Xpub, '/get_gm_child_xpub')
api.add_resource(Gm_Xpub_Verify, '/gm_xpub_verify')
api.add_resource(Gm_Xprv_Sign, '/gm_xprv_sign')
api.add_resource(Get_Gm_P2WPKH_Program, '/get_gm_P2WPKH_program')
api.add_resource(Get_Gm_Address, '/get_gm_address')
api.add_resource(Get_Gm_New_Key, '/get_gm_new_key')
api.add_resource(Get_Gm_New_Address, '/get_gm_new_address')
api.add_resource(Decode_Raw_Tx, '/decode_raw_tx')
api.add_resource(Get_Testnet_Coins, '/get_testnet_coins')
api.add_resource(Get_Gm_Testnet_Coins, '/get_gm_testnet_coins')
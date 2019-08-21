import hashlib
from app.model import key
from app.model import key_gm
from app.model import segwit_addr
import qrcode
import pybase64
from io import BytesIO

# get_path_from_index create xpub path from account key index and current address index
# path: purpose(0x2c=44)/coin_type(btm:0x99)/account_index/change(1 or 0)/address_index
# You can find more details from: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
# You can get more test data from: https://gist.github.com/zcc0721/616eaf337673635fa5c9dd5dbb8dd114
# Please attention:
#   account_index_int >= 1
#   address_index_int >= 1
#   change_bool: true or false
# test data 1:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: true
#   path_list: 2c000000 99000000 01000000 01000000 01000000
# test data 2:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: false
#   path_list: 2c000000 99000000 01000000 00000000 01000000
# test data 3:
#   account_index_int: 3
#   address_index_int: 1
#   change_bool: false
#   path_list: 2c000000 99000000 03000000 00000000 01000000
def get_path_from_index(account_index_int, address_index_int, change_bool):
    path_list = ['2c000000', '99000000']
    account_index_str = (account_index_int).to_bytes(4, byteorder='little').hex()
    path_list.append(account_index_str)
    change_str = '0'
    if change_bool:
        branch_str = (1).to_bytes(4, byteorder='little').hex()
        change_str = '1'
    else:
        branch_str = (0).to_bytes(4, byteorder='little').hex()
    path_list.append(branch_str)
    address_index_str = (address_index_int).to_bytes(4, byteorder='little').hex()
    path_list.append(address_index_str)
    path_str = 'm/44/153/' + str(account_index_int) + '/' + change_str + '/' + str(address_index_int)
    return {
        "path": path_list,
        "path_str": path_str
    }


# create_P2WPKH_program create control program
# You can get more test data from: https://gist.github.com/zcc0721/afa12de04b03b9bfc49985a181ebda80
# Please attention:
#   account_index_int >= 1
#   address_index_int >= 1
#   change_bool: true or false
# test data 1:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: false
#   xpub_str: 3c6664244d2d57168d173c4691dbf8741a67d972b2d3e1b0067eb825e2005d20c5eebd1c26ccad4de5142d7c339bf62cc1fb79a8b3e42a708cd521368dbc9286
#   control_program_str: 0014052620b86a6d5e07311d5019dffa3864ccc8a6bd
# test data 2:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: true
#   xpub_str: 3c6664244d2d57168d173c4691dbf8741a67d972b2d3e1b0067eb825e2005d20c5eebd1c26ccad4de5142d7c339bf62cc1fb79a8b3e42a708cd521368dbc9286
#   control_program: 001478c3aa31753389fcde04d33d0779bdc2840f0ad4
# test data 3:
#   account_index_int: 1
#   address_index_int: 17
#   change_bool: true
#   xpub_str: 3c6664244d2d57168d173c4691dbf8741a67d972b2d3e1b0067eb825e2005d20c5eebd1c26ccad4de5142d7c339bf62cc1fb79a8b3e42a708cd521368dbc9286
#   control_program: 0014eefb8d0688d7960dfbd79bb3aa1bcaa3ec34415d
# test data 4:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: false
#   xpub_str: f744493a021b65814ea149118c98aae8d1e217de29fefb7b2024ca341cd834586ee48bbcf1f4ae801ecb8c6784b044fc62a74c58c816d14537e1573c3e20ce79
#   control_program: 001431f2b90b469e89361225aae370f73e5473b9852b
def create_P2WPKH_program(account_index_int, address_index_int, change_bool, xpub_str):
    path_list = get_path_from_index(account_index_int, address_index_int, change_bool)['path']
    child_xpub_str = key.xpub_to_child_xpub(xpub_str, path_list)['child_xpub']
    child_public_key_str = key.xpub_to_public_key(child_xpub_str)['public_key']
    child_public_key_byte = bytes.fromhex(child_public_key_str)
    
    ripemd160 = hashlib.new('ripemd160')
    ripemd160.update(child_public_key_byte)
    public_key_hash_str = ripemd160.hexdigest()
    control_program_str = '0014' + public_key_hash_str
    return {
        "control_program": control_program_str
    }


# get_gm_P2WPKH_program create control program
# You can get more test data from: https://gist.github.com/zcc0721/5b74b89cae9d75f752a3d2aa9086b9e6
# Please attention:
#   account_index_int >= 1
#   address_index_int >= 1
#   change_bool: true or false
# test data 1:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: false
#   xpub_str: 03c74f3a946940d43e0f8c6da40680c0078e6e1008ca6ea869d57536c31b7ede20adc168c3698fa538fa587c4e519d1eb7a2593f178bfe0c93890a0f09e1634607
#   control_program_str: 001434ca9c0cd289ecb8e899f1bdc0e9470442eae367
# test data 2:
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: true
#   xpub_str: 03c74f3a946940d43e0f8c6da40680c0078e6e1008ca6ea869d57536c31b7ede20adc168c3698fa538fa587c4e519d1eb7a2593f178bfe0c93890a0f09e1634607
#   control_program: 0014ee20c5dfb6730d9abd8bf9c5516b8710bc118271
# test data 3:
#   account_index_int: 1
#   address_index_int: 17
#   change_bool: true
#   xpub_str: 02914d51fcc3b90a87ffe3424995a9db8757a9d67812edd85207c47edc9f7f34e368684ae4d706f68c710fe1dbd20d73a8faaaf34966678a5d58486ac193a851ca
#   control_program: 00149441d8a2e415c63c579d7998563472c1a7c4df2f
# test data 4:
#   account_index_int: 33
#   address_index_int: 44
#   change_bool: true
#   xpub_str: 03603b2eb079257513d253a92ad45ce5ef12cc285dd8c13fc77c95844468f5eb4482f33c577c3a71ac733136b17d68b65a184b225431ab658555735e436fdb13e6
#   control_program: 0014ac49fa87604998fb0ff9d5f83c6c8e803ca84de2
def get_gm_P2WPKH_program(account_index_int, address_index_int, change_bool, xpub_str):
    path_list = get_path_from_index(account_index_int, address_index_int, change_bool)['path']
    child_xpub_str = key_gm.get_gm_child_xpub(xpub_str, path_list)['child_xpub']
    child_public_key_str = key_gm.get_gm_public_key(child_xpub_str)['public_key']
    child_public_key_byte = bytes.fromhex(child_public_key_str)

    ripemd160 = hashlib.new('ripemd160')
    ripemd160.update(child_public_key_byte)
    public_key_hash_str = ripemd160.hexdigest()
    control_program_str = '0014' + public_key_hash_str
    return {
        "control_program": control_program_str
    }


# create_address create address
# You can get more test data from: https://gist.github.com/zcc0721/8f52d0a80a0a9f964e9d9d9a50e940c5
# Please attention:
#   network_str: mainnet/testnet/solonet
# test data 1:
#   control_program_str: 001431f2b90b469e89361225aae370f73e5473b9852b
#   network_str: mainnet
#   address_str: bm1qx8etjz6xn6ynvy394t3hpae723emnpft3nrwej
# test data 2:
#   control_program_str: 0014eefb8d0688d7960dfbd79bb3aa1bcaa3ec34415d
#   network_str: mainnet
#   address_str: bm1qamac6p5g67tqm77hnwe65x7250krgs2avl0nr6
# test data 3:
#   control_program_str: 0014eefb8d0688d7960dfbd79bb3aa1bcaa3ec34415d
#   network_str: testnet
#   address_str: tm1qamac6p5g67tqm77hnwe65x7250krgs2agfwhrt
# test data 4:
#   control_program_str: 0014d234314ea1533dee584417ecb922f904b8dd6c6b
#   network_str: testnet
#   address_str: tm1q6g6rzn4p2v77ukzyzlktjgheqjud6mrt7emxen
# test data 5:
#   control_program_str: 0014eefb8d0688d7960dfbd79bb3aa1bcaa3ec34415d
#   network_str: solonet
#   address_str: sm1qamac6p5g67tqm77hnwe65x7250krgs2adw9jr5
# test data 6:
#   control_program_str: 0014052620b86a6d5e07311d5019dffa3864ccc8a6bd
#   network_str: solonet
#   address_str: sm1qq5nzpwr2d40qwvga2qval73cvnxv3f4aa9xzh9
def create_address(control_program_str, network_str):
    public_key_hash_str = control_program_str[4:]
    if network_str == 'mainnet':
        hrp = 'bm'
    elif network_str == 'testnet':
        hrp = 'tm'
    else:
        hrp = 'sm'
    address_str = segwit_addr.encode(hrp, 0, bytes.fromhex(public_key_hash_str))
    return {
        "address": address_str
    }


# get_gm_address create address
# You can get more test data from: https://gist.github.com/zcc0721/58ff3b33c54616c289dd0b14f75d316c
# Please attention:
#   network_str: gm_testnet/gm_solonet
# test data 1:
#   control_program_str: 0014d234314ea1533dee584417ecb922f904b8dd6c6b
#   network_str: gm_testnet
#   address_str: gm1q6g6rzn4p2v77ukzyzlktjgheqjud6mrtj2c2ef
# test data 2:
#   control_program_str: 0014d234314ea1533dee584417ecb922f904b8dd6c6b
#   network_str: gm_solonet
#   address_str: sm1q6g6rzn4p2v77ukzyzlktjgheqjud6mrtm7srev
# test data 3:
#   control_program_str: 0014eefb8d0688d7960dfbd79bb3aa1bcaa3ec34415d
#   network_str: gm_testnet
#   address_str: gm1qamac6p5g67tqm77hnwe65x7250krgs2ay6dmr3
# test data 4:
#   control_program_str: 0014eefb8d0688d7960dfbd79bb3aa1bcaa3ec34415d
#   network_str: gm_solonet
#   address_str: sm1qamac6p5g67tqm77hnwe65x7250krgs2adw9jr5
def get_gm_address(control_program_str, network_str):
    public_key_hash_str = control_program_str[4:]
    if network_str == 'gm_testnet':
        hrp = 'gm'
    else:
        hrp = 'sm'
    address_str = segwit_addr.encode(hrp, 0, bytes.fromhex(public_key_hash_str))
    return {
        "address": address_str
    }


# create_qrcode_base64 create qrcode, then encode it to base64
# type(s) is str
def create_qrcode_base64(s):
    img = qrcode.make(s)
    buffered = BytesIO()
    img.save(buffered, format="JPEG")
    base64_str = pybase64.b64encode(buffered.getvalue()).decode("utf-8")
    return {
        "base64": base64_str
    }


# create_new_address create address and address qrcode
# test data 1:
#   xpub_str: 8fde12d7c9d6b6cbfbf344edd42f2ed86ae6270b36bab714af5fd5bb3b54adcec4265f1de85ece50f17534e42016ee9404a11fec94ddfadd4a064d27ef3f3f4c
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: False
#   network_str: solonet
#   path: m/44/153/1/0/1
#   control_program: 00147640f3c34fe4b2b298e54e54a4692a47ce47aa5e
#   address: sm1qweq08s60ujet9x89fe22g6f2gl8y02j7lgr5v5
#   address_base64: /9j/4AAQSkZJRgABAQ...
# test data 2:
#   xpub_str: 8fde12d7c9d6b6cbfbf344edd42f2ed86ae6270b36bab714af5fd5bb3b54adcec4265f1de85ece50f17534e42016ee9404a11fec94ddfadd4a064d27ef3f3f4c
#   account_index_int: 12
#   address_index_int: 3
#   change_bool: True
#   network_str: mainnet
#   path: m/44/153/12/1/3
#   control_program: 001458b1477abc46ef81905d25011d36389c0788984b
#   address: bm1qtzc5w74ugmhcryzay5q36d3cnsrc3xztzw6u4y
#   address_base64: /9j/4AAQSkZJRgABAQA...
# test data 3:
#   xpub_str: 8fde12d7c9d6b6cbfbf344edd42f2ed86ae6270b36bab714af5fd5bb3b54adcec4265f1de85ece50f17534e42016ee9404a11fec94ddfadd4a064d27ef3f3f4c
#   account_index_int: 200
#   address_index_int: 1
#   change_bool: True
#   network_str: mainnet
#   path: m/44/153/200/1/1
#   control_program: 00144e5c8757c612c21aa2a0c55f1f8e2ab57cfdefca
#   address: bm1qfewgw47xztpp4g4qc403lr32k470mm724cphhp
#   address_base64: /9j/4AAQSkZJRgABAQA...
def create_new_address(xpub_str, account_index_int, address_index_int, change_bool, network_str):
    path_str = get_path_from_index(account_index_int, address_index_int, change_bool)['path_str']
    control_program_str = create_P2WPKH_program(account_index_int, address_index_int, change_bool, xpub_str)['control_program']
    address_str = create_address(control_program_str, network_str)['address']
    address_base64 = create_qrcode_base64(address_str)['base64']
    return {
        "path": path_str,
        "control_program": control_program_str,
        "address": address_str,
        "address_base64": address_base64
    }


# get_gm_new_address create address and address qrcode
# test data 1:
#   xpub_str: 0269ed53316c1ab5e3b3aa7c1ac3023e250e92a8b1495e1b88fe9564c1c6a49aaeabe1954985189c598214b56563a175fa843e091c9bf941f34a479c8121b196
#   account_index_int: 1
#   address_index_int: 1
#   change_bool: False
#   network_str: solonet
#   path: m/44/153/1/0/1
#   control_program: 00147a9db32f2059c9afc606525a826cd26715364841
#   address: sm1q02wmxteqt8y6l3sx2fdgymxjvu2nvjzpxgal3p
# test data 2:
#   xpub_str: 0254d4f8377fa9cb7839e9525503ff2aa471a82042678082cf45e497362a1dd0e630dee82c0ce53458beafd257852739ab304ca38463ff23f89aa5ee19bc9312
#   account_index_int: 12
#   address_index_int: 3
#   change_bool: True
#   network_str: testnet
#   path: m/44/153/12/1/3
#   control_program: 001443419a20a0ca857c263a871136b64afea3a1cbf3
#   address: gm1qgdqe5g9qe2zhcf36sugnddj2l636rjln0n98yq
# test data 3:
#   xpub_str: 03c8267b52c3918dd8714a30415460b5166f1b3e692b129121787d5bec86e6f57a23f928ec780dcadec36e17d19e46cfb20031c4bc5da26d6d9162f16c2e72da71
#   account_index_int: 200
#   address_index_int: 1
#   change_bool: True
#   network_str: testnet
#   path: m/44/153/200/1/1
#   control_program: 001438b8ad08da405b7565258f2ab94437384d4bdb3c
#   address: gm1q8zu26zx6gpdh2ef93u4tj3ph8px5hkeur9l3x5
def get_gm_new_address(xpub_str, account_index_int, address_index_int, change_bool, network_str):
    path_str = get_path_from_index(account_index_int, address_index_int, change_bool)['path_str']
    control_program_str = create_P2WPKH_program(account_index_int, address_index_int, change_bool, xpub_str)['control_program']
    address_str = get_gm_address(control_program_str, network_str)['address']
    address_base64 = create_qrcode_base64(address_str)['base64']
    return {
        "path": path_str,
        "control_program": control_program_str,
        "address": address_str,
        "address_base64": address_base64
    }

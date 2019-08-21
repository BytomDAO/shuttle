import random
import hashlib
import pbkdf2
import hmac
import ed25519
from app.model.signature import *
from app.model.edwards25519 import *
from app.model.utils import *
from app.model import receiver

# create_key create 128 bits entropy
def create_entropy():
    entropy = random.randint(0, 2**128)
    entropy_str = entropy.to_bytes(16, byteorder='big').hex()
    return {
        "entropy": entropy_str
    }


# entropy_to_mnemonic create mnemonic from 128 bits entropy(the entropy_str length is 32)
# return 12 mnemonics
# You can get more test data from: https://gist.github.com/zcc0721/63aeb5143807950f7b7051fadc08cef0
# test data 1:
#   entropy_str: 1db8b283eb4623e749732a341396e0c9
#   mnemonic_str: buffalo sheriff path story giraffe victory chair grab cross original return napkin
# test data 2:
#   entropy_str: 4d33735a9e92f634d22aecbb4044038d
#   mnemonic_str: essay oppose stove diamond control bounce emerge frown robust acquire abstract brick
# test data 3:
#   entropy_str: 089fe9bf0cac76760bc4b131d938669e
#   mnemonic_str: ancient young hurt bone shuffle deposit congress normal crack six boost despair
def entropy_to_mnemonic(entropy_str):
    mnemonic_str = ""
    mnemonic_length = 12

    # create a 12 elements mnemonic_list 
    mnemonic_list = []
    for _ in range(mnemonic_length):
        mnemonic_list.append('')

    entropy_bytes = bytes.fromhex(entropy_str[:32])
    checksum = hashlib.sha256(entropy_bytes).hexdigest()[:1]
    new_entropy_str = "0" + entropy_str[:32] + checksum
    new_entropy_bytes = bytes.fromhex(new_entropy_str)
    new_entropy_int = int.from_bytes(new_entropy_bytes, byteorder='big')

    file = open('./app/model/english_mnemonic.txt', mode='r')
    word_list = file.readlines()
    file.close()

    for i in range(11, -1, -1):
        word_index = new_entropy_int % 2048
        new_entropy_int = new_entropy_int >> 11
        mnemonic_list[i] = word_list[word_index]

    for i in range(12):
        mnemonic_str += mnemonic_list[i][:-1]
        mnemonic_str += " "
    return {
        "mnemonic": mnemonic_str[:-1]
    }


# mnemonic_to_seed create seed from mnemonic
# You can find more details from: https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki#from-mnemonic-to-seed
# You can get more test data from: https://gist.github.com/zcc0721/4918e891073a9ca6c444ec7490298e82
# test data 1:
#   mnemonic_str: ancient young hurt bone shuffle deposit congress normal crack six boost despair
#   seed_str: afa3a86bbec2f40bb32833fc6324593824c4fc7821ed32eac1f762b5893e56745f66a6c6f2588b3d627680aa4e0e50efd25065097b3daa8c6a19d606838fe7d4
# test data 2:
#   mnemonic_str: rich decrease live pluck friend recipe burden minor similar agent tired horror
#   seed_str: b435f948bd3748ede8f9d6f59728d669939e79c6c885667a5c138e05bbabde1de0dcfcbe0c6112022fbbf0da522f4e224a9c2381016380688b51886248b3156f
# test data 3:
#   mnemonic_str: enough ginger just mutual fit trash loop mule peasant lady market hub
#   seed_str: ecc2bbb6c0492873cdbc81edf56bd896d3b644047879840e357be735b7fa7b6f4af1be7b8d71cc649ac4ca3816f9ccaf11bf49f4effb845f3c19e16eaf8bfcda
def mnemonic_to_seed(mnemonic_str):
    password_str = mnemonic_str
    salt_str = "mnemonic"
    seed_str = pbkdf2.PBKDF2(password_str, salt_str, iterations=2048, digestmodule=hashlib.sha512, macmodule=hmac).hexread(64)
    return {
        "seed": seed_str
    }


# s_str must be >= 32 bytes long and gets rewritten in place.
# This is NOT the same pruning as in Ed25519: it additionally clears the third
# highest bit to ensure subkeys do not overflow the second highest bit.
def prune_root_scalar(s_str):
    s_bytes = bytes.fromhex(s_str)
    s = bytearray(s_bytes)
    s[0] = s[0] & 248
    s[31] = s[31] & 31 # clear top 3 bits
    s[31] = s[31] | 64 # set second highest bit
    return s


# seed_to_root_xprv create rootxprv from seed
# seed_str length is 512 bits.
# root_xprv length is 512 bits.
# You can get more test data from: https://gist.github.com/zcc0721/0aa1b971f4bf93d8f67e25f57b8b97ee
# test data 1:
#   seed_str: afa3a86bbec2f40bb32833fc6324593824c4fc7821ed32eac1f762b5893e56745f66a6c6f2588b3d627680aa4e0e50efd25065097b3daa8c6a19d606838fe7d4
#   root_xprv_str: 302a25c7c0a68a83fa043f594a2db8b44bc871fced553a8a33144b31bc7fb84887c9e75915bb6ba3fd0b9f94a60b7a5897ab9db6a48f888c2559132dba9152b0
# test data 2:
#   seed_str: b435f948bd3748ede8f9d6f59728d669939e79c6c885667a5c138e05bbabde1de0dcfcbe0c6112022fbbf0da522f4e224a9c2381016380688b51886248b3156f
#   root_xprv_str: 6032adeb967ac5ccbf988cf8190817bf9040c8cfd9cdfe3d5e400effb2946946d478b61cc6be936f367ae769eb1dc65c473ee73cac2eb43cf6d5e7c62b7f0062
# test data 3:
#   seed_str: ecc2bbb6c0492873cdbc81edf56bd896d3b644047879840e357be735b7fa7b6f4af1be7b8d71cc649ac4ca3816f9ccaf11bf49f4effb845f3c19e16eaf8bfcda
#   root_xprv_str: a01d6b741b0e74b8d0836ac22b675bbf8e108148ef018d1b000aef1a899a134bd316c0f59e7333520ae1a429504073b2773869e95aa95bb3a4fa0ec76744025c
def seed_to_root_xprv(seed_str):
    hc_str = hmac.HMAC(b'Root', bytes.fromhex(seed_str), digestmod=hashlib.sha512).hexdigest()
    root_xprv_str = prune_root_scalar(hc_str[:64]).hex() + hc_str[64:]
    return {
        "root_xprv": root_xprv_str
    }


# xprv_to_xpub derives new xpub from xprv
# xprv length is 64 bytes.
# xpub length is 64 bytes.
# You can get more test data from: https://gist.github.com/zcc0721/d872a219fa91621d60357278bc62a512
# PLEASE ATTENTION: 
# xprv_bytes = bytes.fromhex(xprv_str)
# xprv_bytes[31] <= 127
# This is the precondition. Please ref: https://github.com/bytom/bytom/blob/dev/crypto/ed25519/internal/edwards25519/edwards25519.go#L958-L963
# test data 1:
#   xprv_str: c003f4bcccf9ad6f05ad2c84fa5ff98430eb8e73de5de232bc29334c7d074759d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   xpub_str: 1b0541a7664cee929edb54d9ef21996b90546918a920a77e1cd6015d97c56563d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
# test data 2:
#   xprv_str: 36667587de27eec684fc4b222276f22a24d9a82e947ee0119148bedd4dec461dd4e1b1d95dfb0f78896677ea1026af7510b41fabd3bd5771311c0cb6968337b2
#   xpub_str: ef0b3a8b0d66523d88f214900101ddb08a2a2a6db28bd8002e5995c1f1cbbc4cd4e1b1d95dfb0f78896677ea1026af7510b41fabd3bd5771311c0cb6968337b2
# test data 3:
#   xprv_str: 74a49c698dbd3c12e36b0b287447d833f74f3937ff132ebff7054baa18623c35a705bb18b82e2ac0384b5127db97016e63609f712bc90e3506cfbea97599f46f
#   xpub_str: 522940d6440fdc45363df2097e9cac29a9a8a33ac339f8b7cff848c199db5a1ca705bb18b82e2ac0384b5127db97016e63609f712bc90e3506cfbea97599f46f
def xprv_to_xpub(xprv_str):
    xprv_bytes = bytes.fromhex(xprv_str)
    scalar = decodeint(xprv_bytes[:len(xprv_bytes)//2])
    buf = encodepoint(scalarmultbase(scalar))
    xpub = buf + xprv_bytes[len(xprv_bytes)//2:]
    xpub_str = xpub.hex()
    return {
        "xpub": xpub_str
    }


# xprv_to_expanded_private_key create expanded private key from xprv
# You can get more test data from: https://gist.github.com/zcc0721/ef0bf2e69f5e92b29d716981f2a8fe7d
# test data 1:
#   xprv_str: 406c82307bf7978d17f3ecfeea7705370e9faef2027affa86c8027c6e11a8a50e231e65bd97048850ae6c39d0f46b63ae70aa24f5aac7877727c430c2201e6d6
#   expanded_private_key_str_xprv: 406c82307bf7978d17f3ecfeea7705370e9faef2027affa86c8027c6e11a8a50d828bf44b1a109c2bbb4c72685858e2f2ab8b405beef1e4ecc12d1ed8511e8eb
# test data 2:
#   xprv_str: 6032adeb967ac5ccbf988cf8190817bf9040c8cfd9cdfe3d5e400effb2946946d478b61cc6be936f367ae769eb1dc65c473ee73cac2eb43cf6d5e7c62b7f0062
#   expanded_private_key_str_xprv: 6032adeb967ac5ccbf988cf8190817bf9040c8cfd9cdfe3d5e400effb2946946ddbb71e7a76595c6bc24937d76085d24315713764cbdf1364ab9091953009cd8
# test data 3:
#   xprv_str: 509a095ad862322641b8d66e84561aae1d4816045167e2c4dfadf464928e114300c0a162d41c0cdf196d61f4492f546e50bfff253b9d5d930d1bb89197cd333d
#   expanded_private_key_str_xprv: 509a095ad862322641b8d66e84561aae1d4816045167e2c4dfadf464928e11432787f5e10f9598f80fb41e4a648b609463c06e625641366f3279658b2b0f5268
def xprv_to_expanded_private_key(xprv_str):
    hc_str = hmac.HMAC(b'Expand', bytes.fromhex(xprv_str), digestmod=hashlib.sha512).hexdigest()
    expanded_private_key_str = xprv_str[:64] + hc_str[64:]
    return {
        "expanded_private_key": expanded_private_key_str
    }


# xpub_to_public_key create 32 bytes public key from xpub
# xpub length is 64 bytes.
# You can get more test data from: https://gist.github.com/zcc0721/9e10f2fa5bd0c8f33aa6dfc87f6aa856
# test data 1:
#   xpub_str: ecc2bbb6c0492873cdbc81edf56bd896d3b644047879840e357be735b7fa7b6f4af1be7b8d71cc649ac4ca3816f9ccaf11bf49f4effb845f3c19e16eaf8bfcda
#   public_key_str: ecc2bbb6c0492873cdbc81edf56bd896d3b644047879840e357be735b7fa7b6f
# test data 2:
#   xpub_str: 406c82307bf7978d17f3ecfeea7705370e9faef2027affa86c8027c6e11a8a50e231e65bd97048850ae6c39d0f46b63ae70aa24f5aac7877727c430c2201e6d6
#   public_key_str: 406c82307bf7978d17f3ecfeea7705370e9faef2027affa86c8027c6e11a8a50
# test data 3:
#   xpub_str: b435f948bd3748ede8f9d6f59728d669939e79c6c885667a5c138e05bbabde1de0dcfcbe0c6112022fbbf0da522f4e224a9c2381016380688b51886248b3156f
#   public_key_str: b435f948bd3748ede8f9d6f59728d669939e79c6c885667a5c138e05bbabde1d
def xpub_to_public_key(xpub_str):
    public_key_str = xpub_str[:64]
    return {
        "public_key": public_key_str
    }


def prune_intermediate_scalar(f):
    f = bytearray(f)
    f[0] = f[0] & 248       # clear bottom 3 bits
    f[29] = f[29] & 1       # clear 7 high bits
    f[30] = 0               # clear 8 bits
    f[31] = 0               # clear 8 bits
    return f


# xprv_to_child_xprv create new xprv through the path
# xprv_str length is 64 bytes.
# path_list item is hex string.
# child_xprv length is 64 bytes.
# You can get more test data from: https://gist.github.com/zcc0721/3377f520954db38070e8e9c80d3a5bfd
# test data 1:
#   xprv_str: 10fdbc41a4d3b8e5a0f50dd3905c1660e7476d4db3dbd9454fa4347500a633531c487e8174ffc0cfa76c3be6833111a9b8cd94446e37a76ee18bb21a7d6ea66b
#   path_list: 010400000000000000
#   path_list: 0100000000000000
#   child_xprv_str: 0813a3abf814e4b4064b9b0492071176d8d98652081aced6fefe2b7363a83353f960274ff5ef195599a765e7bc24eddc2a1e6c73da0e6e0a4b47e65338bea9a6
# test data 2:
#   xprv_str: c003f4bcccf9ad6f05ad2c84fa5ff98430eb8e73de5de232bc29334c7d074759d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   path_list: 00
#   path_list: 00
#   child_xprv_str: b885ac5535c35ae45b51a84b1190f7c31b21acff552c7680413905a9c6084759e9a8f3578fe2973e37d96bad45e8d9255f3b82019f326550d24374aeafece958
# test data 3:
#   xprv_str: 0031615bdf7906a19360f08029354d12eaaedc9046806aefd672e3b93b024e495a95ba63cf47903eb742cd1843a5252118f24c0c496e9213bd42de70f649a798
#   path_list: 00010203
#   child_xprv_str: 20f86339d653bb928ad1f7456279692ac6adf89035f846c6659aaa151c034e497387952cb0dbd6c03bae6742ebe3213b7c8da5805900ab743a653dd3799793eb
# test data 4:
#   xprv_str: 0031615bdf7906a19360f08029354d12eaaedc9046806aefd672e3b93b024e495a95ba63cf47903eb742cd1843a5252118f24c0c496e9213bd42de70f649a798
#   path_list: 00
#   child_xprv_str: 883e65e6e86499bdd170c14d67e62359dd020dd63056a75ff75983a682024e49e8cc52d8e74c5dfd75b0b326c8c97ca7397b7f954ad0b655b8848bfac666f09f
# test data 5:
#   xprv_str: c003f4bcccf9ad6f05ad2c84fa5ff98430eb8e73de5de232bc29334c7d074759d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   path_list: 010203
#   path_list: 7906a1
#   child_xprv_str: 4853a0b00bdcb139e85855d9594e5f641b65218db7c50426946511397e094759bd9de7f2dcad9d7d45389bc94baecaec88aabf58f6e1d832b1f9995a93ec37ea
def xprv_to_child_xprv(xprv_str, path_list):
    for i in range(len(path_list)):
        selector_bytes = bytes.fromhex(path_list[i])
        xpub_str = xprv_to_xpub(xprv_str)['xpub']
        xpub_bytes = bytes.fromhex(xpub_str)
        xprv_bytes = bytes.fromhex(xprv_str)
        hc_bytes = hmac.HMAC(xpub_bytes[32:], b'N'+xpub_bytes[:32]+selector_bytes, digestmod=hashlib.sha512).digest()
        hc_bytearr = bytearray(hc_bytes)
        
        f = hc_bytearr[:32]
        f = prune_intermediate_scalar(f)
        hc_bytearr = f[:32] + hc_bytearr[32:]
        
        carry = 0
        total = 0
        for i in range(32):
            total = xprv_bytes[i] + hc_bytearr[i] + carry
            hc_bytearr[i] = total & 0xff
            carry = total >> 8
        if (total >> 8) != 0:
            print("sum does not fit in 256-bit int")
        xprv_str = hc_bytearr.hex()
        
    child_xprv_str = xprv_str
    return {
        "child_xprv": child_xprv_str
    }


# xpub_to_child_xpub create new xpub through the path
# xpub_str length is 64 bytes.
# path_list item is hex string.
# child_xpub length is 64 bytes.
# You can get more test data from: https://gist.github.com/zcc0721/1dea9eb1edb04f57bc01fecb867301b8
# test data 1:
#   xpub_str: cb22ce197d342d6bb440b0bf13ddd674f367275d28a00f893d7f0b10817690fd01ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   path_list: 010400000000000000
#   path_list: 0100000000000000
#   child_xpub_str: 25405adf9bcefebaa2533631a6bdd5a93108e52ed048c7c49df21a28668768f8d15048473b96fc4d3bc041a881168b41552cabe883221a683aeddc37c1f4644c
# test data 2:
#   xpub_str: cb22ce197d342d6bb440b0bf13ddd674f367275d28a00f893d7f0b10817690fd01ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   path_list: 00
#   path_list: 00
#   child_xpub_str: 1ff4b10aa17eb164a01bedf4f48d55c1bcbd55f28adb85e31c4bad98c070fc4ecb4228fb3f2f848384cc1a9ea82e0b351a551a035dd8ab34e198cfe64df86c79
# test data 3:
#   xpub_str: cb22ce197d342d6bb440b0bf13ddd674f367275d28a00f893d7f0b10817690fd01ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   path_list: 00010203
#   child_xpub_str: 19ab025cd895705c5e2fab8d61e97bcf93670d2c2d6b4cdf06b5347a0cf0527df138d9e540093aad51ed56cf67e6a4b36e6c68327c61593707829339cc9a7f65
# test data 4:
#   xpub_str: ead6415a077b91aa7de32e1cf63350f9351d0298f5accc2cf92ef9429bd1f86c166364ce19322721b7dec84442c3665d97d0e995ba4d01c0f4b19b841379ac90
#   path_list: 00010203
#   path_list: 03ededed
#   path_list: 123456
#   child_xpub_str: c6888c31265519f59975f9fe25a4199735efbb24923648dd880dacb6ed580bdc7b79a9aa09095590175f756c1e11fcb4f8febecb67582c9fea154fd2547cd381
# test data 5:
#   xpub_str: 1b0541a7664cee929edb54d9ef21996b90546918a920a77e1cd6015d97c56563d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   path_list: 010203
#   path_list: 7906a1
#   child_xpub_str: e65c1a9714e2116c6e5d57dee188a53b98dc901a21def5a5ca46fcf78303f4f2bd9de7f2dcad9d7d45389bc94baecaec88aabf58f6e1d832b1f9995a93ec37ea
def xpub_to_child_xpub(xpub_str, path_list):
    for i in range(len(path_list)):
        selector_bytes = bytes.fromhex(path_list[i])
        xpub_bytes = bytes.fromhex(xpub_str)
        hc_bytes = hmac.HMAC(xpub_bytes[32:], b'N'+xpub_bytes[:32]+selector_bytes, digestmod=hashlib.sha512).digest()
        hc_bytearr = bytearray(hc_bytes)

        f = hc_bytearr[:32]
        f = prune_intermediate_scalar(f)
        f = bytes(f)
        scalar = decodeint(f)
        F = scalarmultbase(scalar)

        P = decodepoint(xpub_bytes[:32])
        P = edwards_add(P, F)
        public_key = encodepoint(P)

        xpub_bytes = public_key[:32] + hc_bytes[32:]
        xpub_str = xpub_bytes.hex()

    child_xpub_str = xpub_str
    return {
        "child_xpub": child_xpub_str
    }


# xprv_sign sign message
# xprv_str length is 64 bytes.
# message_str length is variable.
# signature_str length is 64 bytes.
# You can get more test data from: https://gist.github.com/zcc0721/61a26c811a632623678e274cc7e5c10b
# test data 1:
#   xprv_str: c003f4bcccf9ad6f05ad2c84fa5ff98430eb8e73de5de232bc29334c7d074759d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   xpub_str: 1b0541a7664cee929edb54d9ef21996b90546918a920a77e1cd6015d97c56563d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   message_str: a6ce34eec332b32e42ef3407e052d64ac625da6f
#   signature_str: f02f5bb22d8b32f14e88059a786379c26256892f45cf64770c844d0c5de2e52c00307b7bb25fcbb18be13c339a2f511a7c015a8cf81ac681052efe8e50eff00e
# test data 2:
#   xprv_str: 008ce51e3b52ee03eb0ad96c55eb5c9fe8736410518b585a0b7f35b2ab48d24c166364ce19322721b7dec84442c3665d97d0e995ba4d01c0f4b19b841379ac90
#   xpub_str: ead6415a077b91aa7de32e1cf63350f9351d0298f5accc2cf92ef9429bd1f86c166364ce19322721b7dec84442c3665d97d0e995ba4d01c0f4b19b841379ac90
#   message_str: 68656c6c6f206279746f6d      # value is: 'hello bytom'
#   signature_str: 1cc6b0f4031352ffd7a62540f13edddaaebf2df05db7a4926df5513129a8e85dcff1324545a024b16f958239ea67840ced3c2d57bb468dbf0e6cf1d1075f0b0f
# test data 3:
#   xprv_str: 88c0c40fb54ef9c1b90af8cce8dc4c9d54f915074dde93f79ab61cedae03444101ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   xpub_str: cb22ce197d342d6bb440b0bf13ddd674f367275d28a00f893d7f0b10817690fd01ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   message_str: 1246b84985e1ab5f83f4ec2bdf271114666fd3d9e24d12981a3c861b9ed523c6
#   signature_str: ab18f49b23d03295bc2a3f2a7d5bb53a2997bed733e1fc408b50ec834ae7e43f7da40fe5d9d50f6ef2d188e1d27f976aa2586cef1ba00dd098b5c9effa046306
def xprv_sign(xprv_str, message_str):
    xprv_str = xprv_to_expanded_private_key(xprv_str)['expanded_private_key']
    xprv_bytes = bytes.fromhex(xprv_str)
    message_bytes = bytes.fromhex(message_str)
    data_bytes = xprv_bytes[32:64] + message_bytes

    message_digest = hashlib.sha512(data_bytes).digest()
    message_digest = sc_reduce32(message_digest.hex().encode())
    message_digest = bytes.fromhex(message_digest.decode())
    message_digest_reduced = message_digest[0:32]

    scalar = decodeint(message_digest_reduced)
    encoded_r = encodepoint(scalarmultbase(scalar))
    xpub_str = xprv_to_xpub(xprv_str)['xpub']
    xpub_bytes = bytes.fromhex(xpub_str)
    hram_digest_data = encoded_r + xpub_bytes[:32] + message_bytes

    hram_digest = hashlib.sha512(hram_digest_data).digest()
    hram_digest = sc_reduce32(hram_digest.hex().encode())
    hram_digest = bytes.fromhex(hram_digest.decode())
    hram_digest_reduced = hram_digest[0:32]

    sk = xprv_bytes[:32]
    s = sc_muladd(hram_digest_reduced.hex().encode(), sk.hex().encode(), message_digest_reduced.hex().encode())
    s = bytes.fromhex(s.decode())

    signature_bytes = encoded_r + s
    signature_str = signature_bytes.hex()
    return {
        "signature": signature_str
    }


# xpub_verify verify signature
# xpub_str length is 64 bytes.
# message_str length is variable.
# signature_str length is 64 bytes.
# You can get more test data from: https://gist.github.com/zcc0721/61a26c811a632623678e274cc7e5c10b
# test data 1:
#   xprv_str: c003f4bcccf9ad6f05ad2c84fa5ff98430eb8e73de5de232bc29334c7d074759d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   xpub_str: 1b0541a7664cee929edb54d9ef21996b90546918a920a77e1cd6015d97c56563d513bc370335cac51d77f0be5dfe84de024cfee562530b4d873b5f5e2ff4f57c
#   message_str: a6ce34eec332b32e42ef3407e052d64ac625da6f
#   signature_str: f02f5bb22d8b32f14e88059a786379c26256892f45cf64770c844d0c5de2e52c00307b7bb25fcbb18be13c339a2f511a7c015a8cf81ac681052efe8e50eff00e
# test data 2:
#   xprv_str: 008ce51e3b52ee03eb0ad96c55eb5c9fe8736410518b585a0b7f35b2ab48d24c166364ce19322721b7dec84442c3665d97d0e995ba4d01c0f4b19b841379ac90
#   xpub_str: ead6415a077b91aa7de32e1cf63350f9351d0298f5accc2cf92ef9429bd1f86c166364ce19322721b7dec84442c3665d97d0e995ba4d01c0f4b19b841379ac90
#   message_str: 68656c6c6f206279746f6d      # value is: 'hello bytom'
#   signature_str: 1cc6b0f4031352ffd7a62540f13edddaaebf2df05db7a4926df5513129a8e85dcff1324545a024b16f958239ea67840ced3c2d57bb468dbf0e6cf1d1075f0b0f
# test data 3:
#   xprv_str: 88c0c40fb54ef9c1b90af8cce8dc4c9d54f915074dde93f79ab61cedae03444101ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   xpub_str: cb22ce197d342d6bb440b0bf13ddd674f367275d28a00f893d7f0b10817690fd01ff37ac4a07869214c2735bba0175e001abe608db18538e083e1e44430a273b
#   message_str: 1246b84985e1ab5f83f4ec2bdf271114666fd3d9e24d12981a3c861b9ed523c6
#   signature_str: ab18f49b23d03295bc2a3f2a7d5bb53a2997bed733e1fc408b50ec834ae7e43f7da40fe5d9d50f6ef2d188e1d27f976aa2586cef1ba00dd098b5c9effa046306
def xpub_verify(xpub_str, message_str, signature_str):
    result = False
    result = verify(xpub_to_public_key(xpub_str)['public_key'], signature_str, message_str)['result']
    return {
        "result": result
    }


def create_new_key(entropy_str, mnemonic_str):
    if (entropy_str == "") and (mnemonic_str == ""):
        entropy_str = create_entropy()['entropy']
        mnemonic_str = entropy_to_mnemonic(entropy_str)['mnemonic']
    if (entropy_str == "") and (mnemonic_str != ""):
        pass
    if entropy_str != "":
        mnemonic_str = entropy_to_mnemonic(entropy_str)['mnemonic']
    seed_str = mnemonic_to_seed(mnemonic_str)['seed']
    root_xprv_str = seed_to_root_xprv(seed_str)['root_xprv']
    xpub_str = xprv_to_xpub(root_xprv_str)['xpub']
    xprv_base64 = receiver.create_qrcode_base64(root_xprv_str)['base64']
    return {
        "entropy": entropy_str,
        "mnemonic": mnemonic_str,
        "seed": seed_str,
        "xprv": root_xprv_str,
        "xpub": xpub_str,
        "xprv_base64": xprv_base64
    }


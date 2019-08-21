$(function(){
  // 生成新熵
  $('#btnCreateNewKey').click(function(){
    console.log('创建新密钥');
    $.ajax({
      method: 'post',  //get or post
      url: 'https://kit.blockmeta.com/api/v1/get_gm_new_key',
      data: {},
      dataType: 'json',
    }).done(function(data){
      console.log(data);
      layer.msg('创建成功')
      $('#txtEntropy').val(data.entropy)
      $('#txtMnemonics').val(data.mnemonic)
      $('#txtSeed').val(data.seed)
      $('#txtRootXprv').val(data.xprv)
      $('#txtRootXpub').val(data.xpub)
      $('#imgXprvQRCode').attr('src', 'data:image/jpg;base64,' + data.xprv_base64)
    }).fail(function(err){
      layer.alert('创建失败' + err);
    });
  })
  // 清除熵
  $('#btnResetKey').click(function(){
    console.log('清除');
    $('#txtEntropy').val('')
    $('#txtMnemonics').val('')
    $('#txtSeed').val('')
    $('#txtRootXprv').val('')
    $('#txtRootXpub').val('')
    $('#imgXprvQRCode').attr('src', 'data:image/jpg;base64,')
  })
  // 生成新地址
  $('#btnCreateNewAddress').click(function(){
    console.log('创建新地址');
    var xpub_str = $('#txtXpub').val()
    var account_index_int = $('#txtAccountIndex').val()
    var address_index_int = $('#txtAddressIndex').val()
    var change_bool = false
    if ($("input[name='inlineRadioOptionsChange']:checked").val() == "true") {
      change_bool = true
    }
    var network_str = $("input[name='inlineRadioOptionsNetwork']:checked").val()
    $.ajax({
      method: 'post',  //get or post
      url: 'https://kit.blockmeta.com/api/v1/get_gm_new_address',
      data: {
        "xpub_str": xpub_str,
        "account_index_int": account_index_int,
        "address_index_int": address_index_int,
        "change_bool": change_bool,
        "network_str": network_str
      },
      dataType: 'json',
    }).done(function(data){
      console.log(data);
      layer.msg('创建成功')
      $('#txtAddressPath').val(data.path)
      $('#txtControlProgram').val(data.control_program)
      $('#txtAddress').val(data.address)
      $('#imgAddressQRCode').attr('src', 'data:image/jpg;base64,' + data.address_base64)
    }).fail(function(err){
      layer.alert('创建失败' + err);
    });
  })
  // 清除地址
  $('#btnResetAddress').click(function(){
    console.log('清除');
    $('#txtXpub').val('')
    $('#txtAccountIndex').val('1')
    $('#txtAddressIndex').val('1')
    $('#txtAddressPath').val('')
    $('#txtControlProgram').val('')
    $('#txtAddress').val('')
    $('#imgAddressQRCode').attr('src', 'data:image/jpg;base64,')
  })
  // 消息签名
  $('#btnSignMessage').click(function(){
    console.log('创建新签名');
    var xprv_str = $('#txtXprv').val()
    var message_str = $('#txtMessage').val()
    $.ajax({
      method: 'post',  //get or post
      url: 'https://kit.blockmeta.com/api/v1/gm_xprv_sign',
      data: {
        "xprv_str": xprv_str,
        "message_str": message_str
      },
      dataType: 'json',
    }).done(function(data){
      console.log(data);
      layer.msg('签名成功')
      $('#txtSignature').val(data.signature)
    }).fail(function(err){
      layer.alert('签名失败' + err);
    });
  })
  // 清除消息签名
  $('#btnResetSign').click(function(){
    console.log('清除');
    $('#txtXprv').val('')
    $('#txtMessage').val('')
    $('#txtSignature').val('')
  })
  // 验证签名
  $('#btnVerifySignature').click(function(){
    console.log('验证签名');
    var xpub_str = $('#txtXpubVerify').val()
    var message_str = $('#txtMessageVerify').val()
    var signature_str = $('#txtSignatureVerify').val()
    $.ajax({
      method: 'post',  //get or post
      url: 'https://kit.blockmeta.com/api/v1/gm_xpub_verify',
      data: {
        "xpub_str": xpub_str,
        "message_str": message_str,
        "signature_str": signature_str
      },
      dataType: 'json',
    }).done(function(data){
      console.log(data);
      layer.msg('验证成功')
      $('#txtVerifyResult').val(data.result)
    }).fail(function(err){
      layer.alert('验证失败，请检查输入数据格式是否正确' + err);
    });
  })
  // 清除消息签名
  $('#btnResetVerify').click(function(){
    console.log('清除');
    $('#txtXpubVerify').val('')
    $('#txtMessageVerify').val('')
    $('#txtSignatureVerify').val('')
    $('#txtVerifyResult').val('')
  })
  // 发送交易
  $('#btnSubmitTransaction').click(function(){
    console.log('验证签名');
    var raw_transaction_str = $('#txtRawTransaction').val()
    var network_str = $("input[name='inlineRadioOptionsNetworkSubmitTx']:checked").val()
    $.ajax({
      method: 'post',  //get or post
      url: 'https://kit.blockmeta.com/api/v1/submit_transaction',
      data: {
        "raw_transaction_str": raw_transaction_str,
        "network_str": network_str
      },
      dataType: 'json',
    }).done(function(data){
      console.log(data);
      layer.msg('发送成功')
      $('#txtRawTransaction').val('')
    }).fail(function(err){
      layer.alert('发送失败，请检查输入数据格式是否正确' + err);
    });
  })
  // 清除发送交易
  $('#btnResetSubmitTransaction').click(function(){
    console.log('清除');
    $('#txtRawTransaction').val('')
  })
  // 解码原生交易信息
  $('#btnDocodeRawTransaction').click(function(){
    console.log('解码交易');
    var raw_transaction_str = $('#txtRawHexTransaction').val()
    $.ajax({
      method: 'post',  //get or post
      url: 'https://kit.blockmeta.com/api/v1/decode_raw_transaction',
      data: {
        "raw_transaction_str": raw_transaction_str
      },
      dataType: 'json',
    }).done(function(data){
      console.log(data);
      layer.msg('解码成功')
      $('#txtJsonTransaction').val(data.response)
    }).fail(function(err){
      layer.alert('解码失败，请检查输入数据格式是否正确' + err);
    });
  })
  // 清除解码交易
  $('#btnResetDecodeRawTransaction').click(function(){
    console.log('清除');
    $('#txtRawHexTransaction').val('')
    $('#txtJsonTransaction').val('')
  })
});

$(function(){
	window.onload = function() {
		getChain()
		getTransactions()

		setInterval(function(){
			getChain()
			getTransactions()
		},10000)
	}

	function getChain() {
		$.ajax({
			"type": "GET",
			"url": "/chain/list",
			"dataType": "json"
		}).done(function(res){
			drawChainTable(res)
		})
	}

	function getTransactions() {
		$.ajax({
			"type": "GET",
			"url": "/transaction/list",
			"dataType": "json"
		}).done(function(res){
			drawTransactionTable(res)
		})
	}

	function drawChainTable(res) {
		$("#tbl_chain").find("tbody").html("")
		$.each(res, function(k,v){
			let ts = []
			$.each(v.Transactions, function(kk,vv){
				ts.push(vv.From + " -> " + vv.To + " / " + vv.Amount + " Coin")
			})
			let row = [
				"<tr>",
				"<td class='sm col-sm'>",
				v.Index  + ((v.Index == "0" ? "(Genesis)" : "")),
				"</td>",
				"<td class='sm col-sm'>",
				v.Timestamp.replace(" +0900 JST","").replace(/\.\d+$/,""),
				"</td>",
				"<td class='sm col-lg'>",
				ts.join("<br>--------------------<br>"),
				"</td>",
				"<td class='sm col-lg'>",
				v.Hash,
				"</td>",
				"<td class='sm col-lg'>",
				v.PrevHash,
				"</td>",
				"<td class='sm col-lg'>",
				v.Nonce,
				"</td>",
				"<td class='sm col-lg'>",
				v.Miner,
				"</td>",
				"</tr>",
			].join("")
			$(row).prependTo($("#tbl_chain").find("tbody"))
		})
	}


	function drawTransactionTable(res) {
		$("#tbl_transaction").find("tbody").html("")
		$.each(res, function(k,v){
			let ts = []
			let row = [
				"<tr>",
				"<td>",
				v.From,
				"</td>",
				"<td>",
				v.To,
				"</td>",
				"<td>",
				v.Amount,
				"</td>",

				"</tr>",
			].join("")
			$(row).prependTo($("#tbl_transaction").find("tbody"))
		})
	}


	$("#btn_pay").click(function(){
		if (/^[^\s]+/.test($("#From").val()) && /^[^\s]+/.test($("#To").val()) && /^[0-9]+/.test($("#Amount").val())) {
			$.ajax({
				"type": "POST",
				"url": "/transaction/add",
				"dataType": "json",
				"data": JSON.stringify({"From":$("#From").val(), "To": $("#To").val(), "Amount": parseInt($("#Amount").val())}),
				"contentType": "application/json"
			}).done(function(res){
				getChain()
				getTransactions()
				alert(res.msg)
			})
		}
	})


	$("#btn_mining").click(function(){
		if (/^[^\s]+/.test($("#Miner").val())) {
			console.log($("#Miner").val())
			$.ajax({
				"type": "POST",
				"url": "/chain/mining",
				"dataType": "json",
				"data": JSON.stringify({"Miner":$("#Miner").val()}),
				"contentType": "application/json"
			}).done(function(res){
				getChain()
				getTransactions()
				alert(res.msg)
			})
		}
	})

})

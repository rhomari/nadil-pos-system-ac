<template>
  <div
    class="text-center pa-4"
    style="-webkit-user-select: none; -ms-user-select: none; user-select: none"
  >
    <v-dialog :model-value="modelValue" max-width="80mm"  persistent>
      <div class="receipt">
    <div class="receipt-header">
      <h1>{{settings.posname}}</h1>
      <p>{{settings.address}}</p>
      <p>{{settings.phone}}</p>
    </div>
    <span style="font-size:12px; font-weight:bold">Serveur : {{order.waiter}}</span> 
    <table class="receipt-items">
      <thead>
        <tr>
          <th>Produit</th>
          <th >Qté</th>
          <th style="text-align:right">Prix</th>
          <th style="text-align:right">Total</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in JSON.parse(order.content)" :key="item.id" >
          <td>{{item.text}}</td>
          <td>{{item.count}}</td>
          <td style="text-align:right">{{item.price.toFixed(2)}}</td>
          <td style="text-align:right">{{(item.price*item.count).toFixed(2)}}</td>
        </tr>
       
      </tbody>
    </table>

    <div class="receipt-summary">
      
      <span style="float:left; margin-left:50%"><strong>Total:</strong></span><span style="float:right; margin-right:10px"><strong>{{ order.totalticket.toFixed(2)  }}</strong></span> 
    </div>

    <div class="thank-you" style="margin-top:50px;">
      <p style="text-align:left">Merci de votre visite.</p>
      <p style="text-align:left">A bientôt!</p>
     
    </div>
    <div style="color:white; margin-top:10px; text-align:center">
    <v-btn @click="close" style="float:left; margin-right:5px">Fermer</v-btn>
    
    <v-btn @click="print" style="float:right;margin-left:5px;">Imprimer</v-btn>
    </div>
    
  </div>
    </v-dialog>
  </div>
</template>
<script setup>

// eslint-disable-next-line
import {getHostName} from '..\\src\\Utils.js'
const props = defineProps({ modelValue: Boolean, order: Object, settings: Object });
const emit = defineEmits(["close"]);

// eslint-disable-next-line
function close(){
  emit("close")
}
async function print(){
  console.log(props.settings)
  const form = new FormData()
  form.append("number", props.order.number)

  await fetch("https://"+ getHostName()+ "/printreceipt", {method:"post",body:form})
  .then(r=>{
    if (!r.ok) throw Error(r.statusText)
    return r.json()
  })
  .then(data=>{
    console.log(data)
  })
  .catch(err=>{
    console.log(err)
  })
}
</script>

<style>
.receipt {
      
      max-width: 80mm;
      width:80mm;
      margin: 0 auto;
      padding: 10px;
      background: #fff;
      border: 1px solid #ddd;
      border-radius:5px;
      box-shadow: 0 2px 10px rgba(1, 20, 30, 0.1);
      color :black;
    }
    .receipt-header {
      text-align: center;
      margin-bottom: 20px;
    }
    .receipt-header h1 {
      margin: 0;
      font-size: 20px;
    }
    .receipt-header p {
      margin: 0;
      font-size: 12px;
      color: #666;
    }
    .receipt-items {
      width: 100%;
      border-collapse: collapse;
      margin-bottom: 20px;
      font-size:10px;
    }
    .receipt-items th, .receipt-items td {
      text-align: left;
      padding: 8px;
      border-bottom: 1px solid #ddd;
    }
    .receipt-summary {
      text-align: right;
      margin-top: 20px;
    }
    .receipt-summary p {
      margin: 5px 0;
    }
    .thank-you {
      text-align: center;
      margin-top: 20px;
      font-size: 12px;
      color: #555;
    }
</style>

<template>
  <div class="success-alert">
    {{ msg }}
  </div>
</template>

<script>
export default {
  props: {
    msg: {
      type: String,
      required: true
    }
  }
}
</script>

<style scoped>
.success-alert {
  padding: 10px;
  background-color: #d4edda; 
  color: #155724; 
  border: 1px solid #c3e6cb;
  border-radius: 8px;
  margin-top: 15px;
  font-size: 14px;
}
</style>
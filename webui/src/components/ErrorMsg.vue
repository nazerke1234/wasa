<script>
export default {
	props: ['msg']
}
</script>

<template>
	<div class="alert alert-danger" role="alert">
		{{ msg }}
	</div>
</template>

<style>
.custom-error {
  background-color: #fdecea;        /* soft red background */
  border-left: 4px solid #e53935;   /* strong red accent bar */
  color: #b71c1c;                   /* dark red text */
  padding: 12px 16px;
  border-radius: 8px;               /* rounded corners */
  font-size: 14px;
  font-weight: 500;
  margin: 10px 0;
  animation: fadeIn 0.3s ease-in-out;
}

/* Fade-in animation */
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
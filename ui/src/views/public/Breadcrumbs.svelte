<script>
    import { Navigate } from 'svelte-router-spa';
    import { library } from '@fortawesome/fontawesome-svg-core';
    import { faHome as fasHome } from '@fortawesome/free-solid-svg-icons';
    import { FontAwesomeIcon, FontAwesomeLayers, FontAwesomeLayersText } from 'fontawesome-svelte';
</script>

<p><Navigate to=""><FontAwesomeIcon icon={fasHome} size="2x" /></Navigate></p>

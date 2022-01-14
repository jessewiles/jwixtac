<script>
    import { onMount } from 'svelte'
    import { Route } from 'svelte-router-spa';
    import TopHeader from './top_header.svelte';
    import Breadcrumbs from './Breadcrumbs.svelte';
    import Footer from './Footer.svelte';

    export let currentRoute
    const params = {};

</script>

<div class="app">
  <TopHeader />
  <Breadcrumbs />
  <section class="section">
    <Route {currentRoute}  {params} />
  </section>
  <Footer />
</div>

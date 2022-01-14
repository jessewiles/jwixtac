<script>
    import websocketStore from "svelte-websocket-store";
    import { library } from '@fortawesome/fontawesome-svg-core';
    import { faSpinner as fasSpinner } from '@fortawesome/free-solid-svg-icons';
    import { faDownload as fasDownload } from '@fortawesome/free-solid-svg-icons';
    import * as animateScroll from "svelte-scrollto";

    let clicked = 0;

    const initialValue = { };
    export const wsStore = websocketStore("ws://localhost:8088/ws", initialValue, []);

    const feedRoot = "https://hacker-news.firebaseio.com/v0/item";

    let news = {};
    let flips = [];
    let ready = true;

    $: {
        if ($wsStore.message !== undefined) {
            switch($wsStore.channel) {
                case "topstories":
                    $wsStore.message.slice(0, 25).forEach(m => {
                        if (news[m] === undefined) {
                            if (ready) {
                                ready = false;
                                fetch(`${feedRoot}/${m}.json`)
                                .then(res => res.json())
                                .then(jr => {
                                    news[jr.id] = jr;
                                    news = Object.assign({}, news);
                                    ready = true;
                                });
                            }
                        }
                    });
                    break;
                case "imgflip":
                    flips = $wsStore.message.slice(0,25);
                    ready = true;
                    break;
                default:
                    break;
            }

        }
    }

</script>

<div class="all-the-stuff">
    <h2> HACKER NEWS </h2>
    <ul class="hn">
        {#each Object.keys(news) as key}
        <li class="news"><a href={news[key].url} target="hn"> { news[key].title } </a></li>
        {/each}
    </ul>

    <h2> Memes (imgflip) </h2>
    <ul>
        {#each flips as flip}
        <li class="news">
            <div> { flip.name } </div>

            <div>
                <a href={flip.url} target="if"><img src={flip.url} alt={flip.name}/></a>
            </div>
        </li>
        {/each}
    </ul>
</div>



<style>
ul.hn {
    margin: auto;
    width: 65%;
}
h2, li {
    text-align: center;
}
li {
    list-style-type: none;
}
li img {
    max-width: 64px;
    max-height: 64px;
}
  * :global(.myClass) {
    font-style: italic;
  }

  .all-the-stuff {
    padding-bottom: 1in;
  }
</style>

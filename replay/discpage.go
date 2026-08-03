package main

// The discrimination page. Same CSS and same markdown renderer as the
// preference page, because a reader comparing results across the two
// instruments should not also be adjusting to a different-looking screen.
//
// One thing here is not cosmetic. The named voice alternates between pairs, and
// a reader who answers pair 12 against the voice from pair 11 has produced a
// wrong answer for a reason that has nothing to do with the card. So the target
// panel is not a collapsed detail: it sits above the replies, open, and says
// which voice it is every time.

var discHTML = `<!doctype html>
<html lang="en"><head><meta charset="utf-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>cope — blind discrimination</title>
<style>
:root { color-scheme: light dark; --fg:#111; --dim:#666; --bg:#fff; --panel:#f6f6f4; --line:#ddd; --accent:#8a5a00; --good:#2c7a3f; --bad:#a33; }
@media (prefers-color-scheme: dark) {
  :root { --fg:#e6e4e0; --dim:#8f8d88; --bg:#16161a; --panel:#1e1e23; --line:#33333a; --accent:#d0a050; --good:#6ec07f; --bad:#e08080; }
}
* { box-sizing: border-box; }
body { margin:0; background:var(--bg); color:var(--fg); font:16px/1.6 -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,sans-serif; }
header { position:sticky; top:0; background:var(--bg); border-bottom:1px solid var(--line); padding:.7rem 1.2rem; display:flex; gap:1rem; align-items:center; flex-wrap:wrap; z-index:5; }
header h1 { font-size:.95rem; font-weight:600; margin:0; letter-spacing:.02em; }
#bar { flex:1; height:5px; background:var(--panel); border-radius:3px; overflow:hidden; min-width:120px; }
#fill { height:100%; width:0; background:var(--accent); transition:width .2s; }
#count { font-variant-numeric:tabular-nums; color:var(--dim); font-size:.85rem; }
main { max-width:1400px; margin:0 auto; padding:1.2rem; }
.ctx { background:var(--panel); border:1px solid var(--line); border-radius:8px; padding:.9rem 1.1rem; margin-bottom:1rem; font-size:.9rem; }
.ctx summary { cursor:pointer; color:var(--dim); font-size:.82rem; text-transform:uppercase; letter-spacing:.06em; }
.ctx .body { margin-top:.8rem; color:var(--dim); max-height:340px; overflow:auto; }
.target { border:2px solid var(--accent); border-radius:8px; padding:.9rem 1.2rem; margin-bottom:1.2rem; background:var(--panel); }
.target .lbl { font-size:.72rem; text-transform:uppercase; letter-spacing:.08em; color:var(--accent); font-weight:700; margin-bottom:.5rem; }
.target .body { max-height:260px; overflow:auto; font-size:.92rem; }
.prompt { border-left:3px solid var(--accent); padding:.2rem 0 .2rem 1rem; margin-bottom:1.4rem; }
.prompt .lbl { font-size:.72rem; text-transform:uppercase; letter-spacing:.08em; color:var(--dim); }
.cols { display:grid; grid-template-columns:1fr 1fr; gap:1.2rem; }
@media (max-width:900px) { .cols { grid-template-columns:1fr; } }
.card { border:1px solid var(--line); border-radius:8px; padding:1rem 1.2rem; background:var(--panel); }
.card h2 { font-size:.75rem; text-transform:uppercase; letter-spacing:.1em; color:var(--dim); margin:0 0 .8rem; font-weight:600; }
.card .md { overflow-wrap:anywhere; }
.md p { margin:0 0 .8em; } .md ul,.md ol { margin:0 0 .8em; padding-left:1.3em; } .md li { margin:.2em 0; }
.md h1,.md h2,.md h3 { font-size:1em; font-weight:600; margin:1.2em 0 .4em; }
.md code { background:rgba(128,128,128,.16); padding:.1em .35em; border-radius:3px; font-size:.88em; }
.md pre { background:rgba(128,128,128,.14); padding:.7em .9em; border-radius:6px; overflow-x:auto; }
.md pre code { background:none; padding:0; }
.ask { margin:1.6rem auto .4rem; max-width:900px; }
.ask .q { text-align:center; color:var(--dim); font-size:.82rem; margin-bottom:.5rem; }
.pick { display:flex; gap:.7rem; justify-content:center; flex-wrap:wrap; margin-bottom:1.1rem; }
button { font:inherit; padding:.55rem 1.1rem; border:1px solid var(--line); background:var(--bg); color:var(--fg); border-radius:6px; cursor:pointer; }
button:hover { border-color:var(--accent); }
button.on { background:var(--accent); border-color:var(--accent); color:#fff; }
button kbd { font:inherit; font-size:.78em; opacity:.6; margin-right:.4em; }
.hint { text-align:center; color:var(--dim); font-size:.8rem; }
#results { max-width:760px; margin:0 auto; }
table { border-collapse:collapse; width:100%; margin:1rem 0; font-variant-numeric:tabular-nums; }
th,td { text-align:left; padding:.45rem .7rem; border-bottom:1px solid var(--line); font-size:.92rem; }
th { color:var(--dim); font-weight:600; font-size:.78rem; text-transform:uppercase; letter-spacing:.05em; }
.verdict { border-left:3px solid var(--accent); padding:.6rem 1rem; margin:1.2rem 0; }
.good { color:var(--good); } .bad { color:var(--bad); }
</style></head><body>
<header>
  <h1>blind discrimination</h1>
  <div id="bar"><div id="fill"></div></div>
  <span id="count"></span>
  <button id="reveal">reveal</button>
</header>
<main id="app"></main>
<script>
var DATA = JSON.parse(decodeURIComponent(escape(atob("__DATA__"))));
var DESCS = JSON.parse(decodeURIComponent(escape(atob("__DESCS__"))));
var TOTAL = __COUNT__;
var answers = {};
var at = 0;

function esc(s){ var d=document.createElement("div"); d.textContent=s; return d.innerHTML; }

// Enough markdown to read a Claude reply the way the terminal renders one.
function md(src){
  var out=[], lines=src.split("\n"), i=0, para=[], list=null;
  function flushPara(){ if(para.length){ out.push("<p>"+inline(para.join(" "))+"</p>"); para=[]; } }
  function flushList(){ if(list){ out.push("<"+list.tag+">"+list.items.map(function(x){return "<li>"+inline(x)+"</li>";}).join("")+"</"+list.tag+">"); list=null; } }
  function flush(){ flushPara(); flushList(); }
  while(i<lines.length){
    var ln=lines[i];
    if(/^` + "```" + `/.test(ln)){
      flush(); i++; var buf=[];
      while(i<lines.length && !/^` + "```" + `/.test(lines[i])) buf.push(lines[i++]);
      i++; out.push("<pre><code>"+esc(buf.join("\n"))+"</code></pre>"); continue;
    }
    var h=ln.match(/^(#{1,6})\s+(.*)$/);
    if(h){ flush(); out.push("<h3>"+inline(h[2])+"</h3>"); i++; continue; }
    var ul=ln.match(/^\s*[-*+]\s+(.*)$/);
    if(ul){ flushPara(); if(!list||list.tag!=="ul"){ flushList(); list={tag:"ul",items:[]}; } list.items.push(ul[1]); i++; continue; }
    var ol=ln.match(/^\s*\d+[.)]\s+(.*)$/);
    if(ol){ flushPara(); if(!list||list.tag!=="ol"){ flushList(); list={tag:"ol",items:[]}; } list.items.push(ol[1]); i++; continue; }
    if(/^\s*$/.test(ln)){ flush(); i++; continue; }
    flushList(); para.push(ln.trim()); i++;
  }
  flush();
  return out.join("");
}
function inline(s){
  s = esc(s);
  s = s.replace(/` + "`" + `([^` + "`" + `]+)` + "`" + `/g, "<code>$1</code>");
  s = s.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");
  s = s.replace(/(^|[^*])\*([^*]+)\*/g, "$1<em>$2</em>");
  return s;
}

function progress(){
  var n=Object.keys(answers).length;
  document.getElementById("fill").style.width=(100*n/TOTAL)+"%";
  document.getElementById("count").textContent=n+" / "+TOTAL;
}

// The named voice changes between pairs, so it is stated on every screen rather
// than parked in a collapsed panel. A reader answering pair 12 against pair
// 11's voice is a wrong answer the card had nothing to do with.
function render(){
  if(at>=DATA.length){ return results(); }
  var p=DATA[at];
  var app=document.getElementById("app");
  app.innerHTML =
    '<div class="target"><div class="lbl">the voice — which reply was written under it?</div>' +
      '<div class="body md">'+md(DESCS[p.target]||"")+'</div></div>' +
    (p.context ? '<details class="ctx"><summary>conversation so far</summary><div class="body md">'+md(p.context)+'</div></details>' : '') +
    '<div class="prompt"><div class="lbl">the message being answered</div><div class="md">'+md(p.prompt)+'</div></div>' +
    '<div class="cols">' +
      '<div class="card"><h2>reply 1</h2><div class="md">'+md(p.left)+'</div></div>' +
      '<div class="card"><h2>reply 2</h2><div class="md">'+md(p.right)+'</div></div>' +
    '</div>' +
    '<div class="ask"><div class="q">which one was written under the voice above?</div><div class="pick">' +
      '<button data-q="left"><kbd>1</kbd>reply 1</button>' +
      '<button data-q="unsure"><kbd>3</kbd>can\'t tell</button>' +
      '<button data-q="right"><kbd>2</kbd>reply 2</button>' +
    '</div></div>' +
    '<div class="hint">1 / 2 / 3 to answer, → or enter to skip, ← to go back</div>';
  var cur=answers[p.index]||{};
  Array.prototype.forEach.call(app.querySelectorAll("[data-q]"), function(b){
    var v=b.getAttribute("data-q");
    if(cur.q===v) b.className="on";
    b.onclick=function(){ setQ(v); };
  });
  window.scrollTo(0,0);
  progress();
}

function slot(){
  var k=DATA[at].index;
  if(!answers[k]) answers[k]={};
  return answers[k];
}
function setQ(v){ slot().q=v; step(); }
function step(){ at++; render(); }

document.addEventListener("keydown", function(e){
  if(at>=DATA.length) return;
  var q={"1":"left","2":"right","3":"unsure"}[e.key];
  if(q) setQ(q);
  else if(e.key==="ArrowRight"||e.key==="Enter") step();
  else if(e.key==="ArrowLeft" && at>0){ at--; render(); }
});
document.getElementById("reveal").onclick=results;

// Two-sided exact binomial against p=0.5, unsures dropped.
function sign(k,n){
  if(!n) return 1;
  function C(n,k){ var r=1; for(var i=0;i<k;i++) r=r*(n-i)/(i+1); return r; }
  var lo=Math.min(k,n-k), s=0;
  for(var i=0;i<=lo;i++) s+=C(n,i);
  return Math.min(1, 2*s/Math.pow(2,n));
}

// The side the target was actually written under.
function truth(p){ return p.left_cond===p.target ? "left" : "right"; }

// score counts correct calls on one set. pick(p) returns the answer being
// scored, so the reader and the judge go through the same arithmetic.
function score(set, pick){
  var hit=0, miss=0, skip=0;
  DATA.forEach(function(p){
    if(set && p.set!==set) return;
    var v=pick(p);
    if(!v){ return; }
    if(v==="unsure"){ skip++; return; }
    if(v===truth(p)) hit++; else miss++;
  });
  return {hit:hit, miss:miss, skip:skip};
}
function human(p){ var a=answers[p.index]; return a && a.q; }
function judge(p){ return p.judge; }
function hasJudge(){ return DATA.some(function(p){ return p.judge; }); }
function has(set){ return DATA.some(function(p){ return p.set===set; }); }

// The comparator is read out of the data rather than written into the page, so
// a run against any two cards labels its own results.
function against(set){
  for(var i=0;i<DATA.length;i++){
    var p=DATA[i];
    if(p.set!==set) continue;
    return p.left_cond===p.target ? p.right_cond : p.left_cond;
  }
  return "";
}

function row(label, s){
  var n=s.hit+s.miss, p=sign(s.hit,n);
  var rate = n ? (100*s.hit/n).toFixed(0)+"%" : "—";
  var cls = (n && p<0.05) ? (s.hit>s.miss ? "good" : "bad") : "";
  return "<tr><td>"+label+"</td><td>"+s.hit+" / "+n+"</td><td>"+s.skip+"</td><td class=\""+cls+"\">"+rate+
    "</td><td>"+(n ? p.toFixed(3) : "—")+"</td></tr>";
}

function table(caption, pick){
  var rows="";
  ["test","control"].forEach(function(set){
    if(!has(set)) return;
    rows += row(set+" — vs "+against(set), score(set, pick));
  });
  return "<h3>"+caption+"</h3><table><tr><th>set</th><th>correct</th><th>can't tell</th>"+
    "<th>rate</th><th>p</th></tr>" + rows + "</table>";
}

// Agreement is what says whether the judge's rate is about the voice or about
// the model. It counts only pairs both answered with a side.
function agreement(){
  var same=0, n=0;
  DATA.forEach(function(p){
    var h=human(p), j=judge(p);
    if(!h||!j||h==="unsure"||j==="unsure") return;
    n++; if(h===j) same++;
  });
  return {same:same, n:n};
}

function verdict(){
  var t=score("test", human), c=score("control", human);
  var nT=t.hit+t.miss, nC=c.hit+c.miss;
  if(nT<6 && nC<6) return "Too few calls to read.";
  var out="";
  if(nC>=6){
    var pC=sign(c.hit,nC), cr=(100*c.hit/nC).toFixed(0)+"%";
    out += pC<0.05
      ? "The control separates at "+cr+", n="+nC+": the card's voice is visible against no guidance at all, so the instrument can register a hit. "
      : "The control is "+cr+" at n="+nC+" and does not separate. The card cannot be told from bare framing, which is the widest gap on offer — until that moves, nothing below it is readable. ";
  }
  if(nT>=6){
    var pT=sign(t.hit,nT), rate=(100*t.hit/nT).toFixed(0)+"%";
    out += pT<0.05
      ? "The test separates at "+rate+", p="+pT.toFixed(3)+": two deliberate voices, told apart from their descriptions. That is voicing arriving on the page."
      : "No separation on the test: "+rate+", p="+pT.toFixed(3)+". At this n only a large effect is readable, so this rules out the two cards being obviously different and says nothing about a subtle difference.";
  }
  return out;
}

function results(){
  var html = '<div id="results"><h2>results</h2>' + table("you", human);
  if(hasJudge()){
    var a=agreement();
    html += table("the model judge, same pairs", judge);
    html += '<p style="color:var(--dim);font-size:.88rem">You and the judge agree on '+
      a.same+' of '+a.n+' pairs you both called'+(a.n ? ' ('+(100*a.same/a.n).toFixed(0)+'%)' : '')+
      '. Until that agreement is high the judge\'s rate is a fact about a model, not about a voice.</p>';
  }
  html += '<div class="verdict">'+verdict()+'</div>' +
    '<p style="color:var(--dim);font-size:.88rem">p is a two-sided exact binomial against a coin flip, ' +
    '"can\'t tell" dropped. Reload to start over; nothing is saved.</p>' +
    '<p><button id="back">back to the pairs</button></p></div>';
  document.getElementById("app").innerHTML = html;
  document.getElementById("back").onclick=function(){
    if(at>=DATA.length) at=DATA.length-1;
    render();
  };
  progress();
}

render();
</script></body></html>
`

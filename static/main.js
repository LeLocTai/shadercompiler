const inEl = document.getElementById("in");
const outEl = document.getElementById("out");

inEl.value = localStorage.getItem('code') ||
	`struct PSInput
{
	float4 color : COLOR;
};

float4 PSMain(PSInput input) : SV_TARGET
{
	return input.color;
}`

async function compile(code) {
	res = await fetch("/compile", {
		method: "POST",
		body: JSON.stringify({
			code
		})
	});
	return await res.json();
}

async function compileInput() {
	code = inEl.value;
	localStorage.setItem('code', code)

	compiled = await compile(code)
	outEl.value = compiled.Asm;
}

function debounce(func, delay) {
	var timeout;
	return function () {
		var context = this, args = arguments;
		var later = function () {
			timeout = null;
			func.apply(context, args);
		};
		clearTimeout(timeout);
		timeout = setTimeout(later, delay);
	};
};

inEl.onkeyup = debounce(compileInput, 250)

compileInput()
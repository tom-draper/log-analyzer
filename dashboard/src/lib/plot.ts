import type { Config, Layout } from "plotly.js-dist-min";

export const plotConfig: Partial<Config> = {
  responsive: true,
  displayModeBar: false,
  showSendToCloud: false,
};

export function chartLayout(options: { height?: number; title?: string } = {}): Partial<Layout> {
  return {
    title: options.title ? { text: options.title, font: { size: 13, color: "#e4e4e7" }, x: 0 } : undefined,
    hovermode: "x unified",
    plot_bgcolor: "transparent",
    paper_bgcolor: "transparent",
    margin: { t: options.title ? 36 : 12, l: 42, b: 34, r: 12 },
    autosize: true,
    dragmode: false,
    height: options.height ?? 240,
    font: { family: "Geist Variable, ui-sans-serif, system-ui", color: "#a1a1aa", size: 11 },
    hoverlabel: { bgcolor: "#18181b", bordercolor: "#3f3f46", font: { family: "Geist Variable, ui-sans-serif, system-ui", color: "#f4f4f5", size: 12 } },
    yaxis: { gridcolor: "#202023", zerolinecolor: "#27272a", tickfont: { color: "#a1a1aa" }, fixedrange: true },
    xaxis: { gridcolor: "#202023", tickfont: { color: "#a1a1aa" }, fixedrange: true },
  };
}

import "./chartjs.js";

window.chartInstances = window.chartInstances || {};

(function () {
  if (!window.chartScriptInitialized) {
    function getThemeColors() {
      const style = getComputedStyle(document.documentElement);
      return {
        foreground: style.getPropertyValue("--foreground").trim() || "#000",
        background: style.getPropertyValue("--background").trim() || "#fff",
        mutedForeground:
          style.getPropertyValue("--muted-foreground").trim() || "#666",
        border: style.getPropertyValue("--border").trim() || "#ccc",
      };
    }

    function initChart(canvas) {
      if (!canvas || !canvas.id || !canvas.hasAttribute("data-chart-id"))
        return;
      
      if (canvas.hasAttribute("data-initialized")) return;
      canvas.setAttribute("data-initialized", "true");

      if (window.chartInstances[canvas.id]) {
        cleanupChart(canvas);
      }

      const dataId = canvas.getAttribute("data-chart-id");
      const dataElement = document.getElementById(dataId);
      if (!dataElement) return;

      try {
        const chartConfig = JSON.parse(dataElement.textContent);
        const colors = getThemeColors();

        // eslint-disable-next-line no-undef
        Chart.defaults.elements.point.radius = 0;
        // eslint-disable-next-line no-undef
        Chart.defaults.elements.point.hoverRadius = 5;

        const isComplexChart = ["pie", "doughnut", "bar", "radar"].includes(
          chartConfig.type
        );

        const legendOptions = {
          display: chartConfig.showLegend || false,
          labels: { color: colors.foreground },
        };

        const tooltipOptions = {
          backgroundColor: colors.background,
          bodyColor: colors.mutedForeground,
          titleColor: colors.foreground,
          borderColor: colors.border,
          borderWidth: 1,
        };

        const scalesOptions =
          chartConfig.type === "radar"
            ? {
                r: {
                  grid: {
                    color: colors.border,
                    display: chartConfig.showYGrid !== false,
                  },
                  ticks: {
                    color: colors.mutedForeground,
                    backdropColor: "transparent",
                    display: chartConfig.showYLabels !== false,
                  },
                  angleLines: {
                    color: colors.border,
                    display: chartConfig.showXGrid !== false,
                  },
                  pointLabels: {
                    color: colors.foreground,
                    font: { size: 12 },
                  },
                  border: {
                    display: chartConfig.showYAxis !== false,
                    color: colors.border,
                  },
                  beginAtZero: true,
                },
              }
            : {
                x: {
                  beginAtZero: true,
                  display:
                    chartConfig.showXLabels !== false ||
                    chartConfig.showXGrid !== false ||
                    chartConfig.showXAxis !== false,
                  border: {
                    display: chartConfig.showXAxis !== false,
                    color: colors.border,
                  },
                  ticks: {
                    display: chartConfig.showXLabels !== false,
                    color: colors.mutedForeground,
                  },
                  grid: {
                    display: chartConfig.showXGrid !== false,
                    color: colors.border,
                  },
                  stacked: chartConfig.stacked || false,
                },
                y: {
                  offset: true,
                  beginAtZero: true,
                  display:
                    chartConfig.showYLabels !== false ||
                    chartConfig.showYGrid !== false ||
                    chartConfig.showYAxis !== false,
                  border: {
                    display: chartConfig.showYAxis !== false,
                    color: colors.border,
                  },
                  ticks: {
                    display: chartConfig.showYLabels !== false,
                    color: colors.mutedForeground,
                  },
                  grid: {
                    display: chartConfig.showYGrid !== false,
                    color: colors.border,
                  },
                  stacked: chartConfig.stacked || false,
                },
              };

        const finalChartConfig = {
          ...chartConfig,
          options: {
            responsive: true,
            maintainAspectRatio: false,
            interaction: {
              intersect: isComplexChart ? true : false,
              axis: "xy",
              mode: isComplexChart ? "nearest" : "index",
            },
            indexAxis: chartConfig.horizontal ? "y" : "x",
            plugins: {
              legend: legendOptions,
              tooltip: tooltipOptions,
            },
            scales: scalesOptions,
          },
        };

        // eslint-disable-next-line no-undef
        window.chartInstances[canvas.id] = new Chart(canvas, finalChartConfig);
      } catch {}
    }

    function cleanupChart(canvas) {
      if (!canvas || !canvas.id || !window.chartInstances[canvas.id]) return;
      try {
        window.chartInstances[canvas.id].destroy();
      } finally {
        delete window.chartInstances[canvas.id];
      }
    }

    function init(root = document) {
      if (typeof Chart === "undefined") return;

      for (const canvas of root.querySelectorAll("canvas[data-chart-id]:not([data-initialized])")) {
        initChart(canvas);
      }
    }

    function waitForChartAndInit() {
      if (typeof Chart !== "undefined") {
        init();
      } else {
        setTimeout(waitForChartAndInit, 100);
      }
    }

    window.templUI = window.templUI || {};
    window.templUI.chart = { init: init, cleanup: cleanupChart };

    document.addEventListener("DOMContentLoaded", waitForChartAndInit);

    const observer = new MutationObserver(() => {
      let timeout;
      clearTimeout(timeout);
      timeout = setTimeout(() => {
        for (const canvas of document.querySelectorAll(
          "canvas[data-chart-id]"
        )) {
          if (window.chartInstances[canvas.id]) {
            cleanupChart(canvas);
            initChart(canvas);
          }
        }
      }, 50);
    });

    observer.observe(document.documentElement, {
      attributes: true,
      attributeFilter: ["class", "style"],
    });

    window.chartScriptInitialized = true;
  }
})();

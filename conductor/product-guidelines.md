# Product Guidelines - Signal Sentry

## Visual Style: Functional & Technical
Signal Sentry prioritizes high-density data and technical accuracy. The interface should feel like a professional network monitoring tool.
- **Layout:** Use standard ASCII tables for structured data.
- **Color Coding:** Employ standard ANSI colors (Green/Yellow/Red) to denote signal health levels (Excellent/Fair/Poor).
- **Density:** Do not shy away from showing multiple related metrics simultaneously.

## Interaction Mode: Hybrid "Find Best Spot"
The tool's placement optimization mode should balance raw precision with user guidance.
- **Precision:** Display exact dBm changes (e.g., "+2 dBm") to show even minor signal improvements.
- **Encouragement:** Use clear threshold labels (e.g., "MOVED TO EXCELLENT") to indicate significant progress.
- **Proactive Advice:** Provide logic-based suggestions (e.g., "Signal quality is low; try moving away from electronic interference") where possible.

## Error Communication: Technical & Verbose
In the event of an error or connection drop, the tool should provide full context to aid in troubleshooting.
- **Context:** Display HTTP status codes and specific network error messages (e.g., `dial tcp 192.168.12.1:80: i/o timeout`).
- **Transparency:** Never hide the underlying cause of a failure; power users need this information to diagnose gateway or local network issues.

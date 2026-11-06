import { useState } from 'react';

const questions = [
  { label: 'Security strategy', detail: 'A documented, board-aligned security roadmap' },
  { label: 'Incident readiness', detail: 'A tested plan for responding to an incident' },
  { label: 'Risk visibility', detail: 'A current view of your most important risks' },
  { label: 'Leadership reporting', detail: 'Clear security metrics for decision-makers' },
];

export default function ReadinessCalculator() {
  const [selected, setSelected] = useState([]);
  const toggle = (index) => setSelected((current) =>
    current.includes(index) ? current.filter((item) => item !== index) : [...current, index]
  );
  const count = selected.length;

  return (
    <div className="rounded-2xl border border-white/10 bg-[#1c3038] p-6 shadow-2xl shadow-black/10 sm:p-8">
      <div className="flex items-start justify-between gap-4">
        <div>
          <p className="text-xs font-bold uppercase tracking-[0.18em] text-[#c7f36b]">The readiness check</p>
          <h3 className="mt-2 text-xl font-semibold text-white">Where do you stand?</h3>
          <p className="mt-2 text-sm leading-6 text-slate-300">Select the practices your organization has in place.</p>
        </div>
        <span className="rounded-full border border-white/15 px-3 py-1 text-xs text-slate-300">4 practices</span>
      </div>
      <div className="mt-6 space-y-2">
        {questions.map((question, index) => (
          <label key={question.label} className="flex cursor-pointer items-center gap-4 rounded-lg border border-white/10 bg-[#15272e] p-4 transition hover:border-[#c7f36b]/50">
            <input type="checkbox" checked={selected.includes(index)} onChange={() => toggle(index)} className="size-4 accent-[#c7f36b]" />
            <span><span className="block text-sm font-semibold text-white">{question.label}</span><span className="mt-1 block text-xs leading-5 text-slate-400">{question.detail}</span></span>
          </label>
        ))}
      </div>
      <div className="mt-6 border-t border-white/10 pt-5" aria-live="polite">
        <div className="flex items-center justify-between text-sm"><span className="text-slate-300">Practices in place</span><strong className="text-[#c7f36b]">{count} of 4</strong></div>
        <div className="mt-3 h-1.5 overflow-hidden rounded-full bg-white/10"><div className="h-full rounded-full bg-[#c7f36b] transition-all" style={{ width: `${count * 25}%` }} /></div>
        <p className="mt-3 text-xs leading-5 text-slate-400">{count === 4 ? 'A strong foundation. Keep testing and refining it as your business evolves.' : 'Every organization starts somewhere. Let’s identify your next best step.'}</p>
      </div>
    </div>
  );
}

import { ComponentChild } from "preact";
import { EventStatus, FunctionRunStatus } from "../../store/generated";
import classNames from "../../utils/classnames";
import statusStyles from "../../utils/statusStyles";
import { Time } from "../Time";

interface FuncCardProps {
  title: string;
  date: string | number;
  badge?: string | number;
  id: string;
  status: FunctionRunStatus | EventStatus;
  active?: boolean;
  contextualBar?: ComponentChild;
  onClick?: () => void;
}

export default function FuncCard({
  title,
  date,
  badge,
  id,
  status,
  active = false,
  contextualBar,
  onClick,
}: FuncCardProps) {
  const itemStatus = statusStyles(status);

  return (
    <a
      className={classNames(
        active
          ? `outline outline-2 outline-indigo-400 outline-offset-3 bg-slate-900 border-slate-700/50`
          : `hover:bg-slate-800`,
        `px-5 py-3.5 bg-slate-800/50 w-full rounded-lg hover:bg-slate-800/80 block cursor-pointer`
      )}
      onClick={
        onClick
          ? (e) => {
              e.preventDefault();
              onClick();
            }
          : undefined
      }
    >
      <div>
        <div className="flex items-start justify-between">
          <div>
            <span className="text-2xs mt-1 block leading-none">
              <Time date={date} />
            </span>
            <h2 className="text-white mt-2">{title}</h2>
          </div>
          {badge ? (
            <div className="flex items-center px-2 py-2 rounded-sm bg-slate-800 text-2xs leading-none text-slate-50">
              {badge}
            </div>
          ) : null}
        </div>
        <div className="flex items-center justify-between mt-2">
          <span className="text-3xs leading-none">{id}</span>
          <span className="text-3xs leading-none flex items-center">
            <itemStatus.icon />
            <span className="ml-2">{status}</span>
          </span>
        </div>
      </div>

      {contextualBar && (
        <div className="border-t border-slate-700/50 mt-5 pt-3 flex items-center justify-between">
          {contextualBar}
        </div>
      )}
    </a>
  );
}

import React from 'react'

export const AskCodyIcon = React.memo(({ iconColor }: { iconColor: string }) => {
    return (
        <svg width="22" height="22" viewBox="0 0 28 20" fill="none" xmlns="http://www.w3.org/2000/svg">
            <mask id="mask0_3_1189" maskUnits="userSpaceOnUse" x="0" y="0" width="21" height="21">
                <circle cx="10.5" cy="10.5" r="10.5" fill="#D9D9D9" />
            </mask>
            <g mask="url(#mask0_3_1189)">
                <circle cx="10.5" cy="10.5" r="10.5" fill={iconColor} />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M13.5389 4.42041C14.2477 4.42041 14.8223 5.04359 14.8223 5.81233L14.8223 8.28685C14.8223 9.05558 14.2477 9.67876 13.5389 9.67876C12.8301 9.67876 12.2555 9.05558 12.2555 8.28685L12.2555 5.81233C12.2555 5.04359 12.8301 4.42041 13.5389 4.42041Z"
                    fill="white"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M5.41064 7.66828C5.41064 6.89955 5.98524 6.27637 6.69403 6.27637H8.97561C9.68441 6.27637 10.259 6.89955 10.259 7.66828C10.259 8.43702 9.68441 9.0602 8.97561 9.0602H6.69403C5.98524 9.0602 5.41064 8.43702 5.41064 7.66828Z"
                    fill="white"
                />
                <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M16.6062 11.7201C17.063 12.1908 17.0815 12.9739 16.6476 13.4693L16.2438 13.9303C13.0841 17.5376 7.75705 17.4477 4.70229 13.7356C4.28277 13.2258 4.32373 12.4436 4.79378 11.9886C5.26384 11.5336 5.98498 11.5781 6.40451 12.0879C8.57092 14.7205 12.3488 14.7842 14.5896 12.226L14.9934 11.765C15.4274 11.2696 16.1494 11.2495 16.6062 11.7201Z"
                    fill="white"
                />
            </g>
        </svg>
    )
})